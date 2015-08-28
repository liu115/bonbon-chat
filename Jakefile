namespace('build', function () {
    desc('Build binary');
    task('build-binary', [], function () {
        jake.exec('mkdir -p build/src');
        jake.exec('ln -s ../../bonbon build/src/bonbon');
        jake.exec('gom install');
        jake.exec('env GOPATH="`pwd`/build" gom build bonbon/bonbon-server');
    });

    desc('Build everything');
    task('all', ['build:build-binary']);
});

namespace('install', function () {
    desc('Install everything');
    task('all', [], function () {
        var prefix = process.env.prefix;
        if (prefix == undefined)
            prefix = '/usr/local';
        jake.exec("mkdir -p '" + prefix + "/bin'");
        jake.exec("cp bonbon-server '" + prefix + "/bin/bonbon-server'");
    });
});

namespace('clean', function () {
    desc('Clean everything');
    task('all', [], function () {
        // jake.rmRf('_vendor');
        jake.exec('rm -rf bonbon-server');
        jake.exec('rm -rf build');
    });
});

namespace('test', function () {
    desc('Run all tests');
    task('all', []);
});

task('build', ['build:all']);

task('install', ['install:all']);

task('clean', ['clean:all']);

task('test', ['test:all']);

task('default', [], function () {
    jake.Task['clean:all'].invoke();
    jake.Task['build:all'].invoke();
    jake.Task['test:all'].invoke();
});
