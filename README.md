travis-deps
===========

命令行：

travis-deps <TravisDepsConf> 

其中 <TravisDepsConf> 文件格式如下：

{
	token: "<AccessToken>",				# 可选的 AccessToken，对于私有 repo 才需要
	deps: [<Repo1>, <Repo2>, ...]		# 依赖的 Github Repo 列表，比如 ["qiniu/rpc", "qiniu/errors"]
}

其他参数：

* 目标repo：$TRAVIS_REPO_SLUG
* 目标branch：$TRAVIS_BRANCH
* 待测试pr：$TRAVIS_PULL_REQUEST，可选

