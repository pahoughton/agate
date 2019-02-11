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
end
