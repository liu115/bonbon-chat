var NAME = 'bonbon';
var VERSION = '0.1';
var PACKAGE_NAME = NAME + '-' + VERSION;
var TAR_GZ_FILENAME = PACKAGE_NAME + '.tar.gz';

namespace('build', function () {
    desc('Build binary');
    task('build-binary', [], function () {
        jake.exec('mkdir -p build/src;' +
                  '[ -L build/src/bonbon ] || ln -s ../../bonbon build/src/bonbon;' +
                  'gom install;' +
                  'env GOPATH="`pwd`/build" gom build bonbon/bonbon-server');
    });

    desc('Build everything');
    task('all', ['build:build-binary']);
});

namespace('clean', function () {
    desc('Clean everything');
    task('all', [], function () {
        jake.exec('rm -rf vendor');
        jake.exec('rm -rf bonbon-server');
        jake.exec('rm -rf build');
        jake.exec('rm -f deployment/bonbon-*.tar.gz');
        jake.exec("rm -f '" + TAR_GZ_FILENAME + "'");
    });
});

namespace('test', function () {
    // TODO call specific test script

    desc('Run all tests');
    task('all', []);
});

desc('Create tar.gz');
task('package-tar-gz', [], {async: true}, function () { // mark async to avoid task package:docker-image accidentally copying unfinished tar.gz file
    jake.exec("git archive -o '" + TAR_GZ_FILENAME + "' --prefix '" + PACKAGE_NAME + "/' --format tar.gz database", function () {
        complete();
    });
});

desc('Create Docker image');
task('package-docker-image', ['package-tar-gz'], function () {
    jake.exec('cp "' + TAR_GZ_FILENAME + '" deployment/;' +
              'cd deployment; docker build .');
});

task('build', ['build:all']);

task('clean', ['clean:all']);

task('test', ['test:all']);

task('package', ['package:all']);

task('default', [], function () {
    jake.Task['build:all'].invoke();
});
