var amqp = require('amqplib/callback_api');

amqp.connect('amqp://rabbitmq', function (error0, connection) {
    if (error0) {
        throw error0;
    }
    
    connection.createChannel(function (error1, channel) {
        if (error1) {
            throw error1;
        }

        var exchange = 'logs';

        channel.assertExchange(exchange, 'fanout', {
            durable: false
        });

        setInterval(() => {
            var random = Math.random() * 10;
            var msg = process.argv.slice(2).join(' ') || `Weight: ${random.toFixed(2)}`;
            console.log("%s [x] Sent %s", new Date(), msg);
            channel.publish(exchange, '', Buffer.from(msg));
        }, 3000);
    });

//    setTimeout(function () {
//        connection.close();
//        process.exit(0);
//    }, 10 * 3000);
});
