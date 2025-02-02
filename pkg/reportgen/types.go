package reportgen

// Report 定义了报告的基本结构
type Report struct {
	Content   string
	FilePath  string
	Week      string              // 周数
	StartDate string              // 起始日期
	EndDate   string              // 结束日期
	Semester  string              // 学期
	Year      string              // 年份
	Sections  map[string][]string // 各个部分的内容
}

// ReportGenerator 定义了报告生成器的接口
type ReportGenerator interface {
	// Generate 生成报告
	Generate(sourcePath string, params map[string]string) error
	// GetAvailablePeriods 获取可用的时间段
	GetAvailablePeriods(sourcePath string) ([]string, error)
}

// Config 定义了配置结构
type Config struct {
	SourceDir      string
	TargetDir      string
	ReportType     string
	SelectedPeriod string
}

// Section 定义了报告中的各个部分
const (
	TeachingSection      = "教学"
	ListeningSection     = "听课"
	TrainingSection      = "培训学习"
	MiscellaneousSection = "杂事"
)
