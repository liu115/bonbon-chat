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

task('clean', ['clean:all']);

task('test', ['test:all']);

task('default', [], function () {
    jake.Task['clean:all'].invoke();
    jake.Task['build:all'].invoke();
    jake.Task['test:all'].invoke();
});
