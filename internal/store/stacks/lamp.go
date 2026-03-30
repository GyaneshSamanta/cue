package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&LAMPStack{}) }

type LAMPStack struct{}

func (s *LAMPStack) Name() string            { return "lamp" }
func (s *LAMPStack) Description() string     { return "Full LAMP stack: Apache, MySQL, PHP, Composer" }
func (s *LAMPStack) EstimatedSizeMB() int    { return 350 }

func (s *LAMPStack) Components() []store.Component {
	return []store.Component{
		{Name: "Apache", Version: "2.4", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"apache2"}, Darwin: []string{"httpd"}, Windows: []string{"ApacheLounge.Apache"}}},
		{Name: "MySQL", Version: "8.x", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"mysql-server"}, Darwin: []string{"mysql"}, Windows: []string{"Oracle.MySQL"}}},
		{Name: "PHP", Version: "8.2+", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"php", "php-fpm", "php-mysql", "php-curl", "php-mbstring"}, Darwin: []string{"php"}, Windows: []string{"PHP.PHP"}}},
		{Name: "Composer", Version: "2.x", DependsOn: []string{"PHP"},
			InstallMethod: store.InstallMethod{Script: "php -r \"copy('https://getcomposer.org/installer', 'composer-setup.php');\" && php composer-setup.php --install-dir=/usr/local/bin --filename=composer"}},
		{Name: "phpMyAdmin", Version: "latest", Optional: true, OptionalPrompt: "(web-based MySQL admin UI)",
			InstallMethod: store.InstallMethod{Linux: []string{"phpmyadmin"}, Darwin: []string{"phpmyadmin"}}},
	}
}

func (s *LAMPStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Apache", Command: "httpd -v"},
		{Name: "MySQL", Command: "mysql --version"},
		{Name: "PHP", Command: "php -v", Pattern: `PHP 8\.`},
		{Name: "Composer", Command: "composer -V"},
	}
}
