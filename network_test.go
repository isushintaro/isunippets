package isunippets

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

func TestCreateUnixDomainSocket(t *testing.T) {
	assert := assert.New(t)

	// 空きポートを確認
	port, err := GetUnusedPort()
	assert.NoError(err)
	name := fmt.Sprintf("isunippets-%s-%d", t.Name(), port)

	// Unix Domain Socket を作成
	sockPath := path.Join(os.TempDir(), fmt.Sprintf("%s.sock", name))
	listener, err := CreateUnixDomainSocket(sockPath)
	defer listener.Close()
	assert.NoError(err)

	// テスト用のエンドポイントを作成
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	defer e.Close()

	server := httptest.NewUnstartedServer(e.Server.Handler)
	server.Listener = listener
	server.Start()
	defer server.Close()

	// Docker で Nginx を起動
	cmd := exec.Command(
		"docker",
		"run",
		"--name",
		name,
		"-d",
		"--rm",
		"-v",
		"./nginx.conf:/etc/nginx/conf.d/default.conf",
		"-v",
		fmt.Sprintf("%s:/tmp/uds.sock", sockPath),
		"-p",
		fmt.Sprintf("%d:80", port),
		"nginx",
	)
	err = cmd.Run()
	assert.NoError(err)
	defer exec.Command("docker", "stop", name).Run()

	time.Sleep(1 * time.Second)

	// Nginx にアクセス
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://localhost:%d", port))

	res, err := client.R().Get("/")
	assert.NoError(err)
	assert.Equal(200, res.StatusCode())
	assert.Contains(res.String(), "Hello, World!")
}

func TestCreateUnixDomainSocket_SocketExists(t *testing.T) {
	assert := assert.New(t)

	// 空きポートを確認
	port, err := GetUnusedPort()
	assert.NoError(err)
	name := fmt.Sprintf("isunippets-%s-%d", t.Name(), port)

	// Unix Domain Socket を作成
	sockPath := path.Join(os.TempDir(), fmt.Sprintf("%s.sock", name))
	listener, err := CreateUnixDomainSocket(sockPath)
	defer listener.Close()
	assert.NoError(err)

	// 既に存在する Unix Domain Socket を作成
	_, err = CreateUnixDomainSocket(sockPath)
	assert.NoError(err)
}

func TestCreateUnixDomainSocket_FileExists(t *testing.T) {
	assert := assert.New(t)

	// 空きポートを確認
	port, err := GetUnusedPort()
	assert.NoError(err)
	name := fmt.Sprintf("isunippets-%s-%d", t.Name(), port)

	sockPath := path.Join(os.TempDir(), fmt.Sprintf("%s.sock", name))

	// 既に存在するファイルを作成
	file, err := os.Create(sockPath)
	assert.NoError(err)
	defer file.Close()

	_, err = CreateUnixDomainSocket(sockPath)
	assert.Error(err)
}

func TestGetUnusedPort(t *testing.T) {
	assert := assert.New(t)

	port, err := GetUnusedPort()
	assert.NoError(err)
	assert.Greater(port, 0)
}
