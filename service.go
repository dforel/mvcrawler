package mvcrawler

import (
	"github.com/tagDong/mvcrawler/conf"
	"github.com/tagDong/mvcrawler/dhttp"
	"time"
)

var (
	tickDur = 60 * time.Second
)

type Service struct {
	modules map[ModuleType]Module

	analysis   *Analysis
	downloader *Downloader
	hServer    *dhttp.HttpServer
}

func NewService() *Service {
	InitLogger()
	s := new(Service)

	s.initAnalysis()
	s.initDownloader()
	s.initModules()

	//
	s.initHttpServer()

	go s.tick()
	return s
}

func (s *Service) initModules() {
	s.modules = map[ModuleType]Module{}
	for mt, fn := range moduleFunc {
		s.modules[mt] = fn(s.analysis, s.downloader, logger)
	}
}

//初始化分析器
func (s *Service) initAnalysis() {
	config := conf.GetConfig().Common.Analysis
	s.analysis = NewAnalysis(config.ChanSize, config.GoroutineCount)

	logger.Infoln("init analysis ok")
}

//初始化下载器
func (s *Service) initDownloader() {
	config := conf.GetConfig().Common.DownLoad
	s.downloader = NewDownLoader(
		config.OutPath, config.ChanSize, config.GoroutineCount, logger)

	logger.Infoln("init analysis ok")
}

//http服务
func (s *Service) initHttpServer() {
	config := conf.GetConfig()
	s.hServer = dhttp.NewHttpServer(config.Common.HttpAddr)

	//注册路由
	s.hServer.Register("/search", s.search)

	go func() {
		err := s.hServer.Listen()
		if err != nil {
			panic(err)
		}
	}()

	logger.Infof("httpServer start on %s", config.Common.HttpAddr)
}

//定时抓取
func (s *Service) tick() {
	tick := time.NewTicker(tickDur)
	for {
		now := <-tick.C
		logger.Infof("-------- tick %s------", now.String())

	}
}
