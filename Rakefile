# 2019-01-07 (cc) <paul4hough@gmail.com>
#
$runstart = Time.now

at_exit {
  runtime = Time.at(Time.now - $runstart).utc.strftime("%H:%M:%S.%3N")
  puts "run time: #{runtime}"
}


task :default do
  sh 'rake --tasks'
  exit 1
end

task :yamllint do
  sh "yamllint -f parsable .travis.yml .gitlab-ci.yml test config"
end

task :test => [:yamllint] do
  sh 'cd config && go test -mod=vendor'
end

task :build do
  sh 'go build -mod=vendor'
  sh 'cd mock-ticket && go build -mod=vendor'
end

task :build_static do
  require 'git'
  git = Git.open('.')

  branch = git.branch
  commit = git.gcommit('HEAD').sha
  version = File.open('VERSION', &:readline).chomp
  tag = git.tags[-1]

  sh 'go build -mod=vendor ' + \
     "-tags netgo -ldflags '" +\
     "-X main.Version=#{version} " +\
     "-X main.Branch=#{branch} " +\
     "-X main.Revision=#{commit} " +\
     "-X main.BuildDate=#{Time.now.strftime("%Y-%m-%d.%H:%M")} " +\
     "-w -extldflags -static'"
end

task :release => [:test, :build_static] do
  require 'git'
  git = Git.open('.')

  branch = git.branch
  commit = git.gcommit('HEAD').sha
  version = File.open('VERSION', &:readline).chomp
  tag = git.tags[-1]

  if tag.sha != commit
    puts "rev not tagged"
    exit 1
  end
  if tag.name != "v#{version}"
    puts "tag '#{tag.name}' != 'v#{version}' VERSION file "
    exit 1
  end
  puts "branch: #{branch}"
  puts "commit: #{commit}"
  puts "version: #{version}"
  modified = false
  git.status.each do |f|
    if f.type || f.untracked
      mod = f.untracked ? "U" : f.type
      puts "#{mod} " + f.path
      modified = true
    end
  end
  if modified
    puts "modified or untracked files exists"
    exit 1
  end
  sh "mkdir agate-#{version}.amd64"
  sh "cp agate README.md VERSION COPYING agate-#{version}.amd64"
  sh "tar czf agate-#{version}.amd64.tar.gz agate-#{version}.amd64"
  sh "tar tzf agate-#{version}.amd64.tar.gz"
end

task :travis do
  sh "yamllint -f parsable .travis.yml .gitlab-ci.yml test config"
  sh 'cd config && go test'
  sh 'go build'
end
