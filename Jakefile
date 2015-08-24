namespace('build', function () {
    desc('Setup build environments');
    task('setup', [], function () {
        jake.mkdirP('build/src');
        jake.exec('ln -s ../../bonbon build/src/bonbon', function () {
            complete();
        });
    });

    desc('Install dependencies');
    task('install-deps', [], function () {
        jake.exec('gom install', function () {
            complete();
        });
    });

    desc('Compile binary');
    task('compile-binary', [], function () {
        jake.exec('env GOPATH="`pwd`/build" gom build bonbon/bonbon-server', function () {
            complete();
        });
        complete();
    });

    desc('Build everything');
    task('all', ['clean:all', 'build:setup', 'build:install-deps', 'build:compile-binary'], function () {
    });
});

namespace('clean', function () {
    desc('Clean everything');
    task('all', [], function () {
        // jake.rmRf('_vendor');
        jake.rmRf('bonbon-server');
        jake.exec('rm -rf build', function () { // jake.rmRf() follows symlinks and destroy bonbon/ directory. it might be a Jake bug
            complete();
        });
    });
});

namespace('test', function () {
    desc('Run sanity tests');
    task('sanity', [], function () {
        // TODO
    });

    desc('Run sanity tests');
    task('all', ['test:sanity'], function () {
        // TODO
    });
});

task('build', ['build:all'], function () {
});

task('clean', ['clean:all'], function () {
});

task('test', ['test:all'], function () {
});

task('default', ['build:all'], function () {
});
