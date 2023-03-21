package config

type Config struct {
	Mysql MysqlConfig `mapstructure:"mysql" yaml:"mysql"`
	Log   LogConfig   `mapstructure:"log" yaml:"log"`
	//jwt签名密钥
	JwtSignKeyHex string `mapstructure:"jwt_sign_key_hex" yaml:"jwt_sign_key_hex"`
	//jwt加密密钥
	JwtSecretHex string `mapstructure:"jwt_secret_hex" yaml:"jwt_secret_hex"`
	//服务端口号
	ServerPort string `mapstructure:"server_port" yaml:"server_port"`
	//视频配置
	Vedio VedioConfig `mapstructure:"vedio" yaml:"vedio"`
}

type MysqlConfig struct {
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Dbname   string `mapstructure:"dbname" yaml:"dbname"`
	//"10s"
	Timeout string `mapstructure:"timeout" yaml:"timeout"`
	//  设置连接池中空闲连接的最大数量。
	MaxIdleConns int `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	//  设置打开数据库连接的最大数量。
	MaxOpenConns int `mapstructure:"max_open_conns" yaml:"max_open_conns"`
}

type LogConfig struct {
	Path         string `mapstructure:"path" yaml:"path"`
	PanicLogName string `mapstructure:"panic_log_name" yaml:"panic_log_name"`
	//trace,debug,info,warn,error,fatal,panic
	Level string `mapstructure:"level" yaml:"level"`
}

type JwtConfig struct {
	SignKeyHex string `mapstructure:"sign_key_hex" yaml:"sign_key_hex"`
	SecretHex  string `mapstructure:"secret_hex" yaml:"secret_hex"`
}

type VedioConfig struct {
	//视频存储的根目录
	BasePath string `mapstructure:"base_path" yaml:"base_path"`
	// e.g. static
	UrlPrefix string `mapstructure:"url_prefix" yaml:"url_prefix"`
	// e.g. http://localhost:8080
	Domain string `mapstructure:"domain" yaml:"domain"`
}
