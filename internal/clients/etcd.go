package clients

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/galaxy-future/BridgX/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var ErrReviewFailed = errors.New("review key already exists")
var etcdClt *EtcdClient

type EtcdClient struct {
	etcdClient *clientv3.Client
	config     *config.EtcdConfig
}

func NewEtcdClient(config *config.EtcdConfig) (*EtcdClient, error) {
	if etcdClt != nil && equalEtcdConfig(config, etcdClt.config) {
		return etcdClt, nil
	}
	if config != nil {
		etcdClient, err := clientv3.New(clientv3.Config{
			Endpoints:   config.Endpoints,
			DialTimeout: config.DailTimeout,
		})
		if err != nil {
			return nil, err
		}

		err = pingEtcdServer(etcdClient)
		if err != nil {
			return nil, fmt.Errorf("failed to ping etcd :%w", err)
		}
		etcdClt = &EtcdClient{etcdClient: etcdClient}
		return etcdClt, nil
	}
	return nil, errors.New("empty config")
}

func (e *EtcdClient) SyncRun(TTL int, key string, job func() error) error {
	//if not config etcd client , just run this job
	if e.etcdClient == nil {
		return job()
	}
	// lock and run this job
	session, err := concurrency.NewSession(e.etcdClient, concurrency.WithTTL(TTL))
	if err != nil {
		return err
	}
	defer func(session *concurrency.Session) {
		_ = session.Close()
	}(session)

	l := concurrency.NewMutex(session, key)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := l.TryLock(ctx); err != nil {
		return err
	}

	err = job()

	_ = l.Unlock(context.Background())
	return err
}
func pingEtcdServer(etcdClient *clientv3.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := etcdClient.Put(ctx, "ping", "pong")
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdClient) GetConfig(group, dataId string) (string, error) {
	kvs, err := e.etcdClient.KV.Get(context.Background(), fmtKey(group, dataId), clientv3.WithLimit(1))
	if err != nil {
		return "", err
	}
	if len(kvs.Kvs) < 1 {
		return "", nil
	}
	return string(kvs.Kvs[0].Value), nil
}

func fmtKey(group, dataId string) string {
	return group + "/" + dataId
}

func (e *EtcdClient) PublishConfig(group, dataId, content string) error {
	_, err := e.etcdClient.KV.Put(context.Background(), fmtKey(group, dataId), content)
	return err
}

func equalEtcdConfig(conf1, conf2 *config.EtcdConfig) bool {
	if conf1 == conf2 {
		return true
	}
	if reflect.DeepEqual(conf1, conf2) {
		return true
	}
	return false
}
