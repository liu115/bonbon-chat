task default: %w[build]
task build: "bonbon-server"

file "bonbon-server" do
  `gom build bonbon/bonbon-server`
end

file "tester" do
  `gom build test-bonbon/tester`
end

task :run => "bonbon-server" do
  `./bonbon-server`
end

task :test => ["bonbon-server", "tester"] do
  `./bonbon-server &`
  sleep 1
  `./tester`
end

namespace :docker do
  task :build do
    puts `sudo docker build -t dadiuhazon/bonbon .`
  end
end

task :clean do
  `rm bonbon-server tester`
end
