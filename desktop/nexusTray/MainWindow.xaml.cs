using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Runtime.InteropServices;
using System.Text;
using System.Windows;
using System.Windows.Input;
using System.Windows.Interop;
using System.Windows.Threading;

namespace NexusTray
{
    // Classes para Mapear o JSON da API
    public class Contract
    {
        public int Id { get; set; }
        public string Title { get; set; }
        public string CompanyName { get; set; }
        // Propriedade auxiliar para aparecer bonito no ComboBox
        public string DisplayName => $"{CompanyName} - {Title}";
    }

    public class AppointmentHistory
    {
        public int ContractId { get; set; }
        public string Description { get; set; }
    }


    public partial class MainWindow : Window
    {
        private DispatcherTimer _clockTimer;
        private bool _isManualEdit = false;
        // CONFIGURAÇÕES
        private const string API_URL = "http://192.168.0.129:8080/api";
        private const int MY_USER_ID = 1; // <--- SEU ID FIXO AQUI POR ENQUANTO

        // ESTADO
        private DateTime _currentStart;
        private string _lastDesc = "Início do dia";
        private int _lastContractId = 0;
        private bool _isTracking = false;

        private List<AppointmentHistory> _fullHistory = new List<AppointmentHistory>();

        private void CmbContracts_SelectionChanged(object sender, System.Windows.Controls.SelectionChangedEventArgs e)
        {
            if (cmbContracts.SelectedValue is int contractId)
            {
                // Pega todas as descrições únicas usadas para este contrato
                var sugestoes = _fullHistory
                    .Where(x => x.ContractId == contractId)
                    .Select(x => x.Description)
                    .Distinct() // Remove duplicados (vários "Daily")
                    .Reverse()  // As mais recentes primeiro (geralmente vem na ordem cronológica)
                    .Take(20)   // Pega só as ultimas 20 pra não poluir
                    .ToList();

                txtDescription.ItemsSource = sugestoes;
            }
        }

        // HOTKEY (Shift+F10)
        [DllImport("user32.dll")] private static extern bool RegisterHotKey(IntPtr hWnd, int id, uint fsModifiers, uint vk);
        [DllImport("user32.dll")] private static extern bool UnregisterHotKey(IntPtr hWnd, int id);
        private const int HOTKEY_ID = 9000;

        public MainWindow()
        {
            InitializeComponent();
            this.Hide(); // Começa escondido
            _currentStart = DateTime.Now;
            // Configura o relógio para atualizar a UI a cada segundo
            _clockTimer = new DispatcherTimer();
            _clockTimer.Interval = TimeSpan.FromSeconds(1);
            _clockTimer.Tick += (s, args) =>
            {
                // Só atualiza o texto se o usuário NÃO estiver editando manualmente
                if (!_isManualEdit && _isTracking)
                {
                    txtEndTime.Text = DateTime.Now.ToString("HH:mm");
                }
            };
            _clockTimer.Start();

            // Carrega contratos assim que abre
            CarregarContratos();
        }

        protected override void OnSourceInitialized(EventArgs e)
        {
            base.OnSourceInitialized(e);
            IntPtr handle = new WindowInteropHelper(this).Handle;
            HwndSource source = HwndSource.FromHwnd(handle);
            source.AddHook(HwndHook);

            // CORREÇÃO ALT + N:
            // 0x1 = Alt
            // 0x4E = Tecla N (78 em decimal)
            RegisterHotKey(handle, HOTKEY_ID, 0x1, 0x4E);
        }

        private IntPtr HwndHook(IntPtr hwnd, int msg, IntPtr wParam, IntPtr lParam, ref bool handled)
        {
            const int WM_HOTKEY = 0x0312;
            if (msg == WM_HOTKEY && wParam.ToInt32() == HOTKEY_ID)
            {
                AbrirJanela();
                handled = true;
            }
            return IntPtr.Zero;
        }

        private void AbrirJanela()
        {
            // Ao abrir a janela, resetamos o estado
            _isManualEdit = false;

            if (_isTracking)
            {
                lblLastTask.Text = _lastDesc;
                lblStartTime.Text = _currentStart.ToString("HH:mm");

                // Já preenche com a hora atual
                txtEndTime.Text = DateTime.Now.ToString("HH:mm");
                txtEndTime.IsReadOnly = true; // Começa travado pra evitar acidente
                txtEndTime.Opacity = 0.7;
            }
            else
            {
                lblLastTask.Text = "Nenhuma atividade ativa.";
                lblStartTime.Text = "--:--";
                txtEndTime.Text = "--:--";
            }

            txtDescription.Text = "";
            txtDescription.Focus();
            this.Show();
            this.Activate();
        }

        // EVENTO 1: Duplo Clique libera a edição
        private void TxtEndTime_MouseDoubleClick(object sender, System.Windows.Input.MouseButtonEventArgs e)
        {
            if (!_isTracking) return;

            _isManualEdit = true; // Para o relógio automático
            txtEndTime.IsReadOnly = false; // Destrava
            txtEndTime.Opacity = 1.0;
            txtEndTime.SelectAll(); // Seleciona tudo pra facilitar digitar
            txtEndTime.Focus();
        }

        // EVENTO 2: Se clicar para focar (opcional, se quiser liberar sem duplo clique)
        private void TxtEndTime_GotFocus(object sender, RoutedEventArgs e)
        {
            // Se quiser que apenas o Foco já pare o relógio, descomente abaixo:
            // _isManualEdit = true;
        }

        private async void BtnSave_Click(object sender, RoutedEventArgs e)
        {
            if (cmbContracts.SelectedValue == null)
            {
                MessageBox.Show("Selecione um contrato!");
                return;
            }
            if (string.IsNullOrWhiteSpace(txtDescription.Text))
            {
                MessageBox.Show("Escreva o que vai fazer!");
                return;
            }

            int novoContratoId = (int)cmbContracts.SelectedValue;
            string novaDesc = txtDescription.Text;
            DateTime agora = DateTime.Now;
            DateTime fimAnterior = ObterHoraFimCalculada();

            // 1. Se estava rodando algo, fecha e salva no banco
            // LÓGICA CRÍTICA: Definir a hora de fim da anterior
            if (_isTracking)
            {
                // Tenta ler o que está na caixa de texto
                if (DateTime.TryParse(txtEndTime.Text, out DateTime horaManual))
                {
                    // O TryParse pega a data de hoje + a hora digitada.
                    // Se o apontamento virou a noite (começou ontem), isso pode dar bug, 
                    // mas para uso diário simples (mesmo dia), funciona perfeito.
                    fimAnterior = horaManual;
                }
                else
                {
                    // Se digitou bobagem, usa Agora
                    fimAnterior = agora;
                }

                // Salva a anterior com a hora calculada/editada
                await EnviarApontamento(_lastContractId, _lastDesc, _currentStart, fimAnterior);
            }

            _currentStart = fimAnterior;

            _lastContractId = (int)cmbContracts.SelectedValue;
            _lastDesc = txtDescription.Text;
            _isTracking = true;

            MyNotifyIcon.ShowBalloonTip("Nexus", $"Iniciando: {_lastDesc}", Hardcodet.Wpf.TaskbarNotification.BalloonIcon.Info);
            this.Hide();
        }

        private void BtnCancel_Click(object sender, RoutedEventArgs e)
        {
            this.Hide();
        }

        private async void CarregarContratos()
        {
            try
            {
                using (var client = new HttpClient())
                {
                    // A. Carrega Contratos (Isso você já tinha)
                    var jsonContratos = await client.GetStringAsync($"{API_URL}/contracts");
                    var contratos = JsonConvert.DeserializeObject<List<Contract>>(jsonContratos);
                    cmbContracts.ItemsSource = contratos;

                    // B. NOVO: Carrega Histórico de Apontamentos Recentes
                    // (Estou chamando /appointments. Se tiver muitos, o ideal seria uma rota /recent-descriptions)
                    var jsonHistory = await client.GetStringAsync($"{API_URL}/appointments");
                    var appointments = JsonConvert.DeserializeObject<List<AppointmentHistory>>(jsonHistory);

                    if (appointments != null)
                    {
                        _fullHistory = appointments;
                    }
                }
            }
            catch (Exception ex)
            {
                // Silencioso ou Log, para não travar o app se a API cair
                System.Diagnostics.Debug.WriteLine($"Erro ao carregar dados: {ex.Message}");
            }
        }

        private async System.Threading.Tasks.Task EnviarApontamento(int contractId, string desc, DateTime start, DateTime end)
        {
            try
            {
                var payload = new
                {
                    contractId = contractId,
                    userId = MY_USER_ID,
                    startTime = start.ToString("yyyy-MM-ddTHH:mm:ssZ"),
                    endTime = end.ToString("yyyy-MM-ddTHH:mm:ssZ"),
                    description = desc
                };

                var json = JsonConvert.SerializeObject(payload);
                var content = new StringContent(json, Encoding.UTF8, "application/json");

                using (var client = new HttpClient())
                {
                    var response = await client.PostAsync($"{API_URL}/appointments", content);
                    if (!response.IsSuccessStatusCode)
                    {
                        var erro = await response.Content.ReadAsStringAsync();
                        MyNotifyIcon.ShowBalloonTip("Erro Nexus", $"Falha ao salvar: {erro}", Hardcodet.Wpf.TaskbarNotification.BalloonIcon.Error);
                    }
                }
            }
            catch (Exception ex)
            {
                MyNotifyIcon.ShowBalloonTip("Erro Nexus", $"Sem conexão: {ex.Message}", Hardcodet.Wpf.TaskbarNotification.BalloonIcon.Error);
            }
        }

        protected override void OnClosed(EventArgs e)
        {
            // Remove o hook do teclado ao fechar totalmente app
            UnregisterHotKey(new WindowInteropHelper(this).Handle, HOTKEY_ID);
            base.OnClosed(e);
        }

        private void MyNotifyIcon_TrayLeftMouseUp(object sender, RoutedEventArgs e)
        {
            if (this.IsVisible)
            {
                this.Hide();
            }
            else
            {
                AbrirJanela();
            }
        }

        // Botão Vermelho de STOP
        private async void BtnStop_Click(object sender, RoutedEventArgs e)
        {
            if (!_isTracking)
            {
                MessageBox.Show("Não há nenhuma atividade rodando para parar.");
                return;
            }

            // 1. Calcula a hora final (Automática ou Editada Manualmente)
            DateTime fimAtividade = ObterHoraFimCalculada();

            // 2. Salva no banco e espera terminar
            await EnviarApontamento(_lastContractId, _lastDesc, _currentStart, fimAtividade);

            // 3. Reseta o estado para "Parado"
            PararRastreamentoVisual();

            // 4. Esconde a janela e avisa
            this.Hide();
            MyNotifyIcon.ShowBalloonTip("Nexus", "Atividade finalizada. Bom descanso! ☕", Hardcodet.Wpf.TaskbarNotification.BalloonIcon.Info);
        }

        // Função auxiliar para evitar repetir código de hora
        private DateTime ObterHoraFimCalculada()
        {
            if (DateTime.TryParse(txtEndTime.Text, out DateTime horaManual))
            {
                return horaManual; // Usa a hora editada se tiver
            }
            return DateTime.Now; // Senão usa Agora
        }

        // Função que "limpa a casa" visualmente
        private void PararRastreamentoVisual()
        {
            _isTracking = false;
            _isManualEdit = false;

            // Limpa variáveis de memória
            _lastContractId = 0;
            _lastDesc = "";

            // Reseta UI
            lblLastTask.Text = "Nenhuma atividade rodando.";
            lblStartTime.Text = "--:--";
            txtEndTime.Text = "--:--";
            txtEndTime.IsReadOnly = true;
            txtEndTime.Opacity = 0.7;

            // Atualiza ícone da bandeja
            MyNotifyIcon.ToolTipText = "Nexus Tracker (Parado)";
        }

        // Botão direito -> Registrar Horas
        private void MenuRegistrar_Click(object sender, RoutedEventArgs e)
        {
            AbrirJanela();
        }

        // Botão direito -> Sair
        private void MenuSair_Click(object sender, RoutedEventArgs e)
        {
            // Fecha a conexão com o banco ou salva estado se precisar
            // Remove o ícone da bandeja para não ficar "fantasma"
            MyNotifyIcon.Dispose();

            // Encerra a aplicação
            Application.Current.Shutdown();
        }

        private void Window_PreviewKeyDown(object sender, KeyEventArgs e)
        {
            if (e.Key == Key.Escape)
            {
                this.Hide();
                // e.Handled = true; // Opcional: impede que o Esc faça outra coisa se tivesse
            }
        }
    }
}