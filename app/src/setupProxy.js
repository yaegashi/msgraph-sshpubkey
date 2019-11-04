const proxy = require('http-proxy-middleware');
module.exports = function (app) {
    const p = proxy({ target: 'http://localhost:8080', changeOrigin: true });
    app.use('/auth', p);
    app.use('/api', p);
}