package clients

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/config"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func TestSyncRunWithNilConfig(t *testing.T) {
	conf := config.EtcdConfig{
		Endpoints:   []string{"127.0.0.1:2379"},
		DailTimeout: time.Second * 5,
	}
	clt, err := NewEtcdClient(&conf)
	if err != nil {
		t.Error("failed to init sync locker with nil config")
	}

	err = clt.SyncRun(5, "test", func() error {
		time.Sleep(1 * time.Millisecond)
		return nil
	})

	if err != nil {
		t.Errorf("failed to run SyncRun with empty etcd config, err : %v", err)
	}

	tmpErr := fmt.Errorf("tmp errors")
	err = clt.SyncRun(5, "test", func() error {
		return tmpErr
	})

	if err != tmpErr {
		t.Errorf("failed to run SyncRun with empty etcd config, want : %v,  got : %v", tmpErr, err)
	}
}

func TestSyncRunWithEtcdConfig(t *testing.T) {
	conf := config.EtcdConfig{
		Endpoints:   []string{"127.0.0.1:2379"},
		DailTimeout: time.Second * 5,
	}
	clt, err := NewEtcdClient(&conf)
	if err != nil {
		t.Error("failed to init sync locker with nil config")
	}

	count := 0
	err = clt.SyncRun(5, "test", func() error {
		count++
		return nil
	})
	if err != nil {
		t.Errorf("failed to run SyncRun with empty etcd config, err : %v", err)
	}
	if count != 1 {
		t.Errorf("failed to run SyncRun with etcd config, want : %v,  got : %v", 1, count)
	}

	tmpErr := fmt.Errorf("tmp errors")
	err = clt.SyncRun(5, "test", func() error {
		return tmpErr
	})
	if err != tmpErr {
		t.Errorf("failed to run SyncRun with empty etcd config, want : %v,  got : %v", tmpErr, err)
	}

	//并行执行
	count = 0
	var wg sync.WaitGroup
	var err1, err2 error

	t.Log("并行执行 ...")
	for i := 0; i < 10; i++ {
		count = 0
		wg.Add(1)
		go func() {
			defer wg.Done()
			err1 = clt.SyncRun(5, "test", func() error {
				time.Sleep(100 * time.Millisecond)
				count++
				return nil
			})
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			err2 = clt.SyncRun(5, "test", func() error {
				time.Sleep(100 * time.Millisecond)
				count++
				return nil
			})
		}()

		wg.Wait()
		if count != 1 {
			t.Errorf("failed to run SyncRun with etcd config, want : %v,  got : %v", 1, count)
		}
		if (err1 == nil && err2 == nil) || (err1 != concurrency.ErrLocked && err2 != concurrency.ErrLocked) {
			t.Errorf("failed to run SyncRun with etcd config, err1 : %v,  err2 : %v", err1, err2)
		}
	}

	t.Log("串行执行 ...")
	for i := 0; i < 10; i++ {
		count = 0
		//串行执行
		wg.Add(1)
		go func() {
			defer wg.Done()
			err1 = clt.SyncRun(1, "test", func() error {
				time.Sleep(100 * time.Millisecond)
				count++
				return nil
			})
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(200 * time.Millisecond)
			err2 = clt.SyncRun(1, "test", func() error {
				count++
				return nil
			})
		}()

		wg.Wait()
		if count != 2 {
			t.Errorf("failed to run SyncRun with etcd config, want : %v,  got : %v", 2, count)
		}
		if err1 != nil && err2 != nil {
			t.Errorf("failed to run SyncRun with etcd config, err1 : %v,  err2 : %v", err1, err2)
		}
	}

	t.Log("不同的key并行执行 ...")
	for i := 0; i < 10; i++ {
		count = 0
		//串行执行
		wg.Add(1)
		go func() {
			defer wg.Done()
			err1 = clt.SyncRun(1, "test1", func() error {
				time.Sleep(100 * time.Millisecond)
				count++
				return nil
			})
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			err2 = clt.SyncRun(1, "test111", func() error {
				count++
				return nil
			})
		}()

		wg.Wait()
		if count != 2 {
			t.Errorf("failed to run SyncRun with etcd config, want : %v,  got : %v", 2, count)
		}
		if err1 != nil && err2 != nil {
			t.Errorf("failed to run SyncRun with etcd config, err1 : %v,  err2 : %v", err1, err2)
		}
	}

}
