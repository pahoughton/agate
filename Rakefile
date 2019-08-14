# 2019-01-07 (cc) <paul4hough@gmail.com>
#
$runstart = Time.now

at_exit {
  runtime = Time.at(Time.now - $runstart).utc.strftime("%H:%M:%S.%3N")
  puts "run time: #{runtime}"
}

$app = 'agate'
$version = File.open('VERSION', &:readline).chomp
$appver = "#{$app}-#{$version}"


task :default do
  sh 'rake --tasks'
  exit 1
end

desc 'lint yml files'
task :yamllint do
  sh "yamllint -f parsable .travis.yml .gitlab-ci.yml test config"
end

desc 'validate'
task :test, [:name] => [:yamllint] do |tasks, args|
  if args[:name]
    sh "cd #{args[:name]} && go test -mod=vendor -v ./..."
  else
    sh 'go test -mod=vendor -v ./...'
  end
end

desc 'compile'
task :build do
  sh 'go build -mod=vendor'
  sh 'cd mock-ticket && go build -mod=vendor'
  sh 'cd mock-service && go build -mod=vendor'
end

task :vup do
  sh 'cd test && vagrant up'
end
task :vprov do
  sh 'cd test && vagrant provision'
end

desc 'create static binary'
task :build_static do
  require 'git'
  git = Git.open('.')

  branch = git.branch
  commit = git.gcommit('HEAD').sha
  if ENV['CI_COMMIT_TAG']
    version = ENV['CI_COMMIT_TAG']
  else
    version = File.open('VERSION', &:readline).chomp
  end
  sh 'go build -mod=vendor ' + \
     "-tags netgo -ldflags '" +\
     "-X main.Version=#{version} " +\
     "-X main.Branch=#{branch} " +\
     "-X main.Revision=#{commit} " +\
     "-X main.BuildDate=#{Time.now.strftime("%Y-%m-%d.%H:%M")} " +\
     "-w -extldflags -static'"
end

desc 'create release tarball'
task :release => [:test, :build_static] do
  require 'git'
  git = Git.open('.')

  branch = git.branch
  commit = git.gcommit('HEAD').sha
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
  puts "version: #{$appver}"
  sh "test -d #{$appver}.amd64 || mkdir #{$appver}.amd64"
  sh "cp #{$app} README.md VERSION COPYING #{$appver}.amd64"
  sh "tar cvzf #{$appver}.amd64.tar.gz #{$appver}.amd64"
end

desc "create #{$appver}.amd64.tar.gz"
task :tarball => [:build_static] do
  puts "version: #{$appver}"
  sh "test -d #{$appver}.amd64 || mkdir #{$appver}.amd64"
  sh "cp #{$app} README.md VERSION COPYING #{$appver}.amd64"
  sh "tar cvzf #{$appver}.amd64.tar.gz #{$appver}.amd64"
end

desc 'tavis validation'
task :travis do
  sh "yamllint -f parsable .travis.yml .gitlab-ci.yml test config"
  sh 'go test -v ./...'
end
