package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		expectCfg *Config
		expectErr error
	}{
		{
			name: "GetConfigOK",
			yaml: `
logger:
  logLevel: DEBUG
  logFile: ./out.log
  logMode:
  rewriteLog: true
database:
  port: 8080
  host: localhost
  dbName: mydb
  user: myuser
  password: mypass`,
			expectCfg: &Config{
				Logger: Logger{
					LogLevel:   "DEBUG",
					LogFile:    "./out.log",
					RewriteLog: true,
				},
				Database: Database{
					Port:     8080,
					Host:     "localhost",
					User:     "myuser",
					Password: "mypass",
					DbName:   "mydb",
				},
			},
			expectErr: nil,
		},
		{
			name: "GetConfigErrorFileExtraSpace",
			yaml: `
logger:
    logLevel: DEBUG
  logFileEnable: true
  logStdoutEnable: true
  maxSize: 500
  maxAge: 28
  maxBackups: 7
  rewriteLog: true`,
			expectCfg: nil,
			expectErr: errors.New("While parsing config: yaml: line 1: did not find expected key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			path := "./config.yaml"

			file, err := os.Create(path)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			defer os.Remove(path)

			file.WriteString(tt.yaml)
			defer file.Close()

			cfg, err := GetConfig()

			if err != nil && tt.expectErr == nil {
				t.Errorf("unexpected error %v, expect nil", err)

			} else if err != nil && tt.expectErr != nil {

				gotErr := strings.Split(fmt.Sprintf("%v", err), "err: ")[1]
				expErr := fmt.Sprintf("%v", tt.expectErr)
				if gotErr != expErr {
					t.Errorf("unexpected error %v,\nexpect %v", gotErr, expErr)
				}
			}

			if !reflect.DeepEqual(cfg, tt.expectCfg) {
				t.Errorf("\nexpect\n%v, \ngot \n%v", tt.expectCfg, cfg)
			}
		})
	}
}

/* ТЕСТИРОВАТЬ ОТДЕЛЬНО
func TestReadFlags(t *testing.T) {
	cfg := &Config{}
	readFlags(cfg)
	tests := []struct {
		name string
		args []string
		want Config
	}{
		{
			name: "all flags",
			args: []string{
				"--loglevel=DEBUG",
				"--logFile=./myapp.log",
				"--logMode=",
				"--rewriteLog=true",
				"--port=8080",
				"--host=localhost",
				"--dbName=mydb",
				"--user=myuser",
				"--password=mypassword",
				"--mqttLogin=myloginmqtt",
				"--mqttPassword=mypasswordmqtt",
				"--mqttHost=localhost",
				"--mqttPort=8080",
				"--mqttDomainName=test",
				"--mqttClientId=test",
				"--metricUrl=http://localhost:9998/metrics",
				"--sleepService=10",
				"--sleepFail=10",
				"--configPath=./",
			},
			want: Config{
				Logger: Logger{
					LogLevel:   "DEBUG",
					LogFile:    "./myapp.log",
					RewriteLog: true,
				},
				Database: Database{
					Port:     8080,
					Host:     "localhost",
					DbName:   "mydb",
					User:     "myuser",
					Password: "mypassword",
				},
				MqttConnect: MqttConnect{
					MqttLogin:      "myloginmqtt",
					MqttPassword:   "mypasswordmqtt",
					MqttHost:       "localhost",
					MqttPort:       "8080",
					MqttDomainName: "test",
					MqttClientId:   "test",
				},
				Metric: Metric{
					MetricUrl: "http://localhost:9998/metrics",
				},
				Sleeps: Sleeps{
					SleepService: 10,
					SleepFail:    10,
				},
			},
		},
	}
	// Run tests
	for _, tt := range tests {
		// Set command-line arguments
		err := flag.CommandLine.Parse(tt.args)
		if err != nil {
			t.Fatalf("test %s: %v", tt.name, err)
		}
		// Compare the resulting Config struct with the expected one
		if *cfg != tt.want {
			t.Errorf("test %s readFlags(%v)\n got: %v\nwant: %v", tt.name, tt.args, *cfg, tt.want)
			continue
		}
		t.Log("Good test", tt.name, *cfg)
	}
}
*/
