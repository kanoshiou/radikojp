package tui

import (
	"fmt"
	"strings"
	"time"

	"radikojp/api"
	"radikojp/config"
	"radikojp/hook"
	"radikojp/model"
	"radikojp/player"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FocusMode 焦点模式
type FocusMode int

const (
	FocusStations FocusMode = iota
	FocusRegion
)

// KeyMap 快捷键
type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Select    key.Binding
	VolUp     key.Binding
	VolDown   key.Binding
	Mute      key.Binding
	Reconnect key.Binding
	Quit      key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.VolUp, k.VolDown, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Select},
		{k.VolUp, k.VolDown, k.Mute, k.Reconnect, k.Quit},
	}
}

var DefaultKeyMap = KeyMap{
	Up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑", "上移")),
	Down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓", "下移")),
	Left:      key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←", "左")),
	Right:     key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→", "右")),
	Select:    key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("Enter", "选择")),
	VolUp:     key.NewBinding(key.WithKeys("+", "="), key.WithHelp("+", "音量+")),
	VolDown:   key.NewBinding(key.WithKeys("-", "_"), key.WithHelp("-", "音量-")),
	Mute:      key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "静音")),
	Reconnect: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "重连")),
	Quit:      key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("Esc", "退出/返回")),
}

// 样式
var (
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#10B981")
	accentColor    = lipgloss.Color("#F59E0B")
	textColor      = lipgloss.Color("#CDD6F4")
	dimTextColor   = lipgloss.Color("#6C7086")
	playingColor   = lipgloss.Color("#A6E3A1")
	regionColor    = lipgloss.Color("#89B4FA")
	warningColor   = lipgloss.Color("#FAB387")

	titleStyle              = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	regionItemStyle         = lipgloss.NewStyle().Foreground(textColor)
	regionSelectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E1E2E")).Background(regionColor).Bold(true).Padding(0, 1)
	regionCurrentStyle      = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)
	stationNameStyle        = lipgloss.NewStyle().Foreground(textColor)
	stationIDStyle          = lipgloss.NewStyle().Foreground(dimTextColor)
	stationSelectedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E1E2E")).Background(primaryColor).Bold(true).Padding(0, 1)
	stationPlayingStyle     = lipgloss.NewStyle().Foreground(playingColor).Bold(true)
	stationSelectedPlayingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E1E2E")).Background(secondaryColor).Bold(true).Padding(0, 1)
	statusStyle             = lipgloss.NewStyle().Foreground(dimTextColor)
	errorStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("#F38BA8"))
	volumeStyle             = lipgloss.NewStyle().Foreground(accentColor)
	focusIndicatorStyle     = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	programStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#CBA6F7"))
	nowPlayingStyle         = lipgloss.NewStyle().Foreground(playingColor).Bold(true)
	reconnectStyle          = lipgloss.NewStyle().Foreground(warningColor)
)

// PlayingInfo 保存正在播放的电台信息
type PlayingInfo struct {
	StationID      string
	StationName    string
	CurrentProgram string
}

// SharedState 共享状态
type SharedState struct {
	Player        *player.FFmpegPlayer
	AuthToken     string
	Volume        float64
	Muted         bool
	CurrentAreaID string
	Playing       *PlayingInfo
}

// Model TUI 模型
type Model struct {
	stations      []model.Station
	cursor        int
	width         int
	height        int
	keys          KeyMap
	statusMessage string
	errorMessage  string
	shared        *SharedState
	autoPlay      bool
	autoPlayIdx   int

	areas        []model.Area
	currentArea  int
	selectedArea int
	isLoading    bool
	focus        FocusMode
}

// 消息类型
type autoPlayMsg struct{}
type stationsLoadedMsg struct {
	stations []model.Station
	err      error
}
type playResultMsg struct {
	err         error
	stationIdx  int
	stationID   string
	stationName string
}
type reconnectResultMsg struct{ err error }
type programUpdateMsg struct{ program string }
type tickMsg struct{}

func NewModel(stations []model.Station, authToken string, initialVolume float64, lastStationID string, areaID string) Model {
	areas := model.AllAreas()

	currentAreaIdx := 0
	for i, area := range areas {
		if area.ID == areaID {
			currentAreaIdx = i
			break
		}
	}

	defaultIdx := 0
	autoPlayIdx := -1
	for i, s := range stations {
		if s.ID == lastStationID {
			defaultIdx = i
			autoPlayIdx = i
			break
		}
	}

	if autoPlayIdx == -1 {
		for i, s := range stations {
			if s.ID == "QRR" {
				defaultIdx = i
				autoPlayIdx = i
				break
			}
		}
	}

	if autoPlayIdx == -1 && len(stations) > 0 {
		autoPlayIdx = 0
	}

	p := player.NewFFmpegPlayer(authToken, initialVolume)

	shared := &SharedState{
		Player:        p,
		AuthToken:     authToken,
		Volume:        initialVolume,
		Muted:         false,
		CurrentAreaID: areaID,
		Playing:       nil,
	}

	p.SetReconnectCallback(func() string {
		return hook.Auth(shared.CurrentAreaID)
	})

	return Model{
		stations:      stations,
		cursor:        defaultIdx,
		keys:          DefaultKeyMap,
		statusMessage: "",
		shared:        shared,
		autoPlay:      true,
		autoPlayIdx:   autoPlayIdx,
		areas:         areas,
		currentArea:   currentAreaIdx,
		selectedArea:  currentAreaIdx,
		focus:         FocusStations,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return autoPlayMsg{} },
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func fetchProgramCmd(stationID string) tea.Cmd {
	return func() tea.Msg {
		prog, err := api.GetCurrentProgram(stationID)
		if err != nil || prog == nil {
			return programUpdateMsg{program: ""}
		}
		return programUpdateMsg{program: prog.Title}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		// 检查重连状态
		if m.shared.Player != nil {
			status := m.shared.Player.GetReconnectStatus()
			if status == player.ReconnectSuccess {
				m.statusMessage = "重连成功"
				m.shared.Player.ClearReconnectStatus()
			} else if status == player.ReconnectFailed {
				m.errorMessage = "重连失败: " + m.shared.Player.GetLastError()
				m.shared.Player.ClearReconnectStatus()
			}
		}
		
		// 每30秒刷新一次节目信息
		var cmd tea.Cmd
		if m.shared.Playing != nil && time.Now().Second()%30 == 0 {
			cmd = fetchProgramCmd(m.shared.Playing.StationID)
		}
		return m, tea.Batch(cmd, tickCmd())

	case programUpdateMsg:
		if m.shared.Playing != nil {
			m.shared.Playing.CurrentProgram = msg.program
		}
		return m, nil

	case autoPlayMsg:
		if m.autoPlay && m.autoPlayIdx >= 0 && m.autoPlayIdx < len(m.stations) {
			m.autoPlay = false
			m.cursor = m.autoPlayIdx
			return m, m.playStation()
		}
		return m, nil

	case stationsLoadedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.errorMessage = fmt.Sprintf("加载失败: %v", msg.err)
		} else {
			m.stations = msg.stations
			m.shared.CurrentAreaID = m.getCurrentAreaID()
			m.cursor = 0
			m.statusMessage = fmt.Sprintf("已切换到 %s", m.getCurrentAreaName())
			m.saveAreaConfig()
		}
		return m, nil

	case playResultMsg:
		if msg.err != nil {
			m.errorMessage = fmt.Sprintf("播放失败: %v", msg.err)
			m.statusMessage = ""
		} else {
			m.shared.Playing = &PlayingInfo{
				StationID:   msg.stationID,
				StationName: msg.stationName,
			}
			m.statusMessage = ""
			m.errorMessage = ""
			m.saveConfig()
			return m, fetchProgramCmd(msg.stationID)
		}
		return m, nil

	case reconnectResultMsg:
		if msg.err != nil {
			m.errorMessage = fmt.Sprintf("重连失败: %v", msg.err)
		} else {
			m.statusMessage = "重连成功"
		}
		return m, nil

	case tea.KeyMsg:
		if m.isLoading {
			return m, nil
		}
		m.errorMessage = ""
		m.statusMessage = ""

		if m.focus == FocusRegion {
			return m.handleRegionKeys(msg)
		}
		return m.handleStationKeys(msg)
	}

	return m, nil
}

func (m Model) handleStationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.focus = FocusRegion
			m.selectedArea = m.currentArea
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if m.cursor < len(m.stations)-1 {
			m.cursor++
		}
		return m, nil

	case key.Matches(msg, m.keys.Left):
		if m.currentArea > 0 {
			m.currentArea--
			m.selectedArea = m.currentArea
			return m, m.loadStationsForCurrentArea()
		}
		return m, nil

	case key.Matches(msg, m.keys.Right):
		if m.currentArea < len(m.areas)-1 {
			m.currentArea++
			m.selectedArea = m.currentArea
			return m, m.loadStationsForCurrentArea()
		}
		return m, nil

	case key.Matches(msg, m.keys.Select):
		return m, m.playStation()

	case key.Matches(msg, m.keys.VolUp):
		if m.shared.Player != nil {
			m.shared.Player.IncreaseVolume(0.05)
			m.shared.Volume = m.shared.Player.GetVolume()
			m.shared.Muted = false
			m.saveConfig()
		}
		return m, nil

	case key.Matches(msg, m.keys.VolDown):
		if m.shared.Player != nil {
			m.shared.Player.DecreaseVolume(0.05)
			m.shared.Volume = m.shared.Player.GetVolume()
			m.shared.Muted = false
			m.saveConfig()
		}
		return m, nil

	case key.Matches(msg, m.keys.Mute):
		if m.shared.Player != nil {
			m.shared.Player.ToggleMute()
			m.shared.Muted = m.shared.Player.IsMuted()
		}
		return m, nil

	case key.Matches(msg, m.keys.Reconnect):
		if m.shared.Player != nil && m.shared.Playing != nil {
			return m, m.reconnect()
		}
		return m, nil

	case key.Matches(msg, m.keys.Quit):
		m.saveConfig()
		if m.shared.Player != nil {
			m.shared.Player.Stop()
		}
		return m, tea.Quit

	case msg.String() >= "0" && msg.String() <= "9":
		if m.shared.Player != nil {
			vol := float64(msg.String()[0]-'0') / 10.0
			m.shared.Player.SetVolume(vol)
			m.shared.Volume = vol
			m.shared.Muted = false
			m.saveConfig()
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleRegionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left):
		if m.selectedArea > 0 {
			m.selectedArea--
		}
		return m, nil

	case key.Matches(msg, m.keys.Right):
		if m.selectedArea < len(m.areas)-1 {
			m.selectedArea++
		}
		return m, nil

	case key.Matches(msg, m.keys.Down), key.Matches(msg, m.keys.Quit):
		m.focus = FocusStations
		m.selectedArea = m.currentArea
		return m, nil

	case key.Matches(msg, m.keys.Select):
		if m.selectedArea != m.currentArea {
			m.currentArea = m.selectedArea
			m.focus = FocusStations
			return m, m.loadStationsForCurrentArea()
		}
		m.focus = FocusStations
		return m, nil
	}
	return m, nil
}

func (m *Model) getCurrentAreaID() string {
	if m.currentArea >= 0 && m.currentArea < len(m.areas) {
		return m.areas[m.currentArea].ID
	}
	return "JP13"
}

func (m *Model) getCurrentAreaName() string {
	if m.currentArea >= 0 && m.currentArea < len(m.areas) {
		return m.areas[m.currentArea].Name
	}
	return "東京"
}

func (m *Model) loadStationsForCurrentArea() tea.Cmd {
	m.isLoading = true
	m.statusMessage = fmt.Sprintf("加载 %s ...", m.getCurrentAreaName())
	areaID := m.getCurrentAreaID()
	return func() tea.Msg {
		stations, err := api.GetStations(areaID)
		return stationsLoadedMsg{stations: stations, err: err}
	}
}

func (m *Model) saveConfig() {
	if m.shared.Playing != nil {
		volume := m.shared.Volume
		if m.shared.Player != nil {
			volume = m.shared.Player.GetVolume()
		}
		go config.SaveConfig(m.shared.Playing.StationID, volume, m.getCurrentAreaID())
	}
}

func (m *Model) saveAreaConfig() {
	volume := m.shared.Volume
	if m.shared.Player != nil {
		volume = m.shared.Player.GetVolume()
	}
	stationID := ""
	if m.shared.Playing != nil {
		stationID = m.shared.Playing.StationID
	}
	go config.SaveConfig(stationID, volume, m.getCurrentAreaID())
}

func (m *Model) playStation() tea.Cmd {
	stationIdx := m.cursor
	station := m.stations[stationIdx]
	shared := m.shared

	return func() tea.Msg {
		playlistURLs, err := api.GetStreamURLs(station.ID)
		if err != nil {
			return playResultMsg{err: err, stationIdx: stationIdx}
		}
		if len(playlistURLs) == 0 {
			return playResultMsg{err: fmt.Errorf("无可用流"), stationIdx: stationIdx}
		}

		lsid := "5e586af5ccb3b0b2498abfb19eaa8472"
		lastUrl := playlistURLs[len(playlistURLs)-1]
		finalStreamUrl := fmt.Sprintf("%s?station_id=%s&l=30&lsid=%s&type=b", lastUrl, station.ID, lsid)

		shared.Player.Stop()
		time.Sleep(100 * time.Millisecond)
		
		// 使用已有的 token，不重新获取
		err = shared.Player.Play(finalStreamUrl)
		return playResultMsg{
			err:         err,
			stationIdx:  stationIdx,
			stationID:   station.ID,
			stationName: station.Name,
		}
	}
}

func (m *Model) reconnect() tea.Cmd {
	shared := m.shared
	return func() tea.Msg {
		if shared.Player != nil {
			return reconnectResultMsg{err: shared.Player.Reconnect()}
		}
		return reconnectResultMsg{err: fmt.Errorf("播放器未初始化")}
	}
}

// View 渲染 - 固定底部布局
func (m Model) View() string {
	// 计算可用高度
	totalHeight := m.height
	if totalHeight == 0 {
		totalHeight = 24 // 默认高度
	}

	// 固定区域高度
	headerHeight := 3  // 标题 + 地区 + 分隔线
	footerHeight := 3  // 分隔线 + 播放信息 + 帮助
	
	// 内容区域高度
	contentHeight := totalHeight - headerHeight - footerHeight
	if contentHeight < 5 {
		contentHeight = 5
	}

	var content strings.Builder

	// === 头部 ===
	// 标题 + 音量
	title := titleStyle.Render("📻 Radiko")
	volBar := m.renderVolume()
	content.WriteString(fmt.Sprintf("%s  %s\n", title, volBar))

	// 地区行
	content.WriteString(m.renderRegionLine() + "\n")
	content.WriteString(strings.Repeat("─", 50) + "\n")

	// === 内容区域 ===
	contentLines := m.renderContent(contentHeight)
	content.WriteString(contentLines)

	// 填充空行使底部固定
	currentLines := strings.Count(content.String(), "\n")
	targetLines := totalHeight - footerHeight
	for i := currentLines; i < targetLines; i++ {
		content.WriteString("\n")
	}

	// === 底部固定区域 ===
	content.WriteString(m.renderFooter())

	return content.String()
}

// renderContent 渲染内容区域
func (m Model) renderContent(maxHeight int) string {
	var lines []string

	if m.isLoading {
		lines = append(lines, fmt.Sprintf("⏳ %s", m.statusMessage))
		return strings.Join(lines, "\n") + "\n"
	}

	// 电台列表
	maxVisible := maxHeight - 2 // 留出状态消息空间
	if maxVisible > len(m.stations) {
		maxVisible = len(m.stations)
	}
	if maxVisible < 3 {
		maxVisible = 3
	}

	startIdx := 0
	if m.cursor >= maxVisible {
		startIdx = m.cursor - maxVisible + 1
	}
	endIdx := startIdx + maxVisible
	if endIdx > len(m.stations) {
		endIdx = len(m.stations)
		startIdx = endIdx - maxVisible
		if startIdx < 0 {
			startIdx = 0
		}
	}

	if startIdx > 0 {
		lines = append(lines, statusStyle.Render("  ↑ 更多"))
	}

	for i := startIdx; i < endIdx; i++ {
		station := m.stations[i]
		isSelected := i == m.cursor && m.focus == FocusStations
		isPlaying := m.shared.Playing != nil && m.shared.Playing.StationID == station.ID

		prefix := "  "
		if isPlaying {
			prefix = "▶ "
		}

		var styled string
		switch {
		case isSelected && isPlaying:
			text := fmt.Sprintf("%s%s %s", prefix, station.Name, station.ID)
			styled = stationSelectedPlayingStyle.Render(text)
		case isSelected:
			text := fmt.Sprintf("%s%s %s", prefix, station.Name, station.ID)
			styled = stationSelectedStyle.Render(text)
		case isPlaying:
			styled = stationPlayingStyle.Render(prefix+station.Name) + " " + stationIDStyle.Render(station.ID)
		default:
			styled = stationNameStyle.Render(prefix+station.Name) + " " + stationIDStyle.Render(station.ID)
		}
		lines = append(lines, styled)
	}

	if endIdx < len(m.stations) {
		lines = append(lines, statusStyle.Render("  ↓ 更多"))
	}

	// 状态/错误消息
	if m.errorMessage != "" {
		lines = append(lines, errorStyle.Render("✗ "+m.errorMessage))
	} else if m.statusMessage != "" {
		lines = append(lines, statusStyle.Render(m.statusMessage))
	}

	return strings.Join(lines, "\n") + "\n"
}

// renderFooter 渲染固定底部
func (m Model) renderFooter() string {
	var lines []string
	
	lines = append(lines, strings.Repeat("─", 50))

	// 播放信息 + 重连状态
	var playLine string
	if m.shared.Playing != nil {
		playLine = nowPlayingStyle.Render("▶ ") + m.shared.Playing.StationName + " " + stationIDStyle.Render(m.shared.Playing.StationID)
		if m.shared.Playing.CurrentProgram != "" {
			playLine += "  " + programStyle.Render("♪ "+m.shared.Playing.CurrentProgram)
		}
		
		// 检查重连状态
		if m.shared.Player != nil {
			switch m.shared.Player.GetReconnectStatus() {
			case player.ReconnectStarted:
				playLine += "  " + reconnectStyle.Render("🔄 重连中...")
			case player.ReconnectAuth:
				playLine += "  " + reconnectStyle.Render("🔑 获取认证...")
			case player.ReconnectPlaying:
				playLine += "  " + reconnectStyle.Render("▶ 恢复播放...")
			}
		}
	} else {
		playLine = statusStyle.Render("未播放")
	}
	lines = append(lines, playLine)

	// 帮助
	if m.focus == FocusRegion {
		lines = append(lines, statusStyle.Render("← → 选择  Enter 确认  ↓/Esc 返回"))
	} else {
		lines = append(lines, statusStyle.Render("↑↓ 选择  Enter 播放  ←→ 切地区  +- 音量  m 静音  r 重连  Esc 退出"))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderVolume() string {
	vol := int(m.shared.Volume * 100)
	if m.shared.Player != nil {
		vol = int(m.shared.Player.GetVolume() * 100)
	}
	if m.shared.Muted {
		return statusStyle.Render(fmt.Sprintf("🔇 %d%%", vol))
	}
	return volumeStyle.Render(fmt.Sprintf("🔊 %d%%", vol))
}

func (m Model) renderRegionLine() string {
	var parts []string

	if m.focus == FocusRegion {
		parts = append(parts, focusIndicatorStyle.Render("▶ "))
	} else {
		parts = append(parts, "  ")
	}

	visibleCount := 5
	startIdx := m.selectedArea - visibleCount/2
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := startIdx + visibleCount
	if endIdx > len(m.areas) {
		endIdx = len(m.areas)
		startIdx = endIdx - visibleCount
		if startIdx < 0 {
			startIdx = 0
		}
	}

	if startIdx > 0 {
		parts = append(parts, statusStyle.Render("◀ "))
	}

	for i := startIdx; i < endIdx; i++ {
		area := m.areas[i]
		var styled string
		if m.focus == FocusRegion && i == m.selectedArea {
			styled = regionSelectedStyle.Render(area.Name)
		} else if i == m.currentArea {
			styled = regionCurrentStyle.Render(area.Name)
		} else {
			styled = regionItemStyle.Render(area.Name)
		}
		parts = append(parts, styled)
		if i < endIdx-1 {
			parts = append(parts, " ")
		}
	}

	if endIdx < len(m.areas) {
		parts = append(parts, statusStyle.Render(" ▶"))
	}

	parts = append(parts, statusStyle.Render(fmt.Sprintf(" [%d/%d]", m.selectedArea+1, len(m.areas))))
	return strings.Join(parts, "")
}

func Run(stations []model.Station, authToken string, cfg config.Config) error {
	m := NewModel(stations, authToken, cfg.Volume, cfg.LastStationID, cfg.AreaID)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()

	if m.shared.Player != nil {
		m.shared.Player.Stop()
	}
	return err
}
