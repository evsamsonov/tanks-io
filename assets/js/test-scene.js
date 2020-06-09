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

        this.socket.on('newPlayer', (playerInfo) => {
            console.log(playerInfo);
        });

        this.socket.on('currentPlayers', (players) => {
            players = JSON.parse(players);

            players.forEach((player) => {
                if (player.hasOwnProperty('id') && player.id === this.socket.id) {
                    this.addPlayer(player);
                }
            });
        });
    }

    update() {}

    addPlayer(playerInfo) {
        this.tank = this.add.image(playerInfo.x, playerInfo.y, 'tank').setOrigin(0.5, 0.5);
    }
}
