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

task :build do
  sh 'go build'
end
