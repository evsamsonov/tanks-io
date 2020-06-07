"use strict";

class TestScene extends Phaser.Scene
{
    preload() {
        this.load.spritesheet('tank', 'img/tank.png', { frameWidth: 60, frameHeight: 60 });
    }

    create() {
        this.socket = io();

        this.socket.on('connect', function () {
            console.log(arguments);
        });

        this.socket.on('connect_error', function () {
            console.log(arguments);
        });

        this.socket.on('error', function () {
            console.log(arguments);
        });

        this.socket.on('connect_timeout', function () {
            console.log(arguments);
        });

        this.socket.on('disconnect', () => {
            console.log(arguments);
        })
    }

    update() {}
}
