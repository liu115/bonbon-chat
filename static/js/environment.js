// 設定開發時與佈屬時的差異

module.exports = function () {
	if (process.env.silent == true || process.env.NODE_ENV == 'production') {
		console.log = function () {};
	}
}