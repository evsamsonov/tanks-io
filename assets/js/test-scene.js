"use strict";

class TestScene extends Phaser.Scene
{
    preload() {
        this.load.spritesheet('tank', 'img/tank.png', { frameWidth: 60, frameHeight: 60 });
    }

    create() {
        this.socket = io();
        this.otherPlayers = this.physics.add.group();

        this.socket.on('newPlayer', (playerInfo) => {
            this.addOtherPlayer(JSON.parse(playerInfo));
        });

        this.socket.on('currentPlayers', (players) => {
            players = JSON.parse(players);

            players.forEach((player) => {
                if (player.id === this.socket.id) {
                    this.addPlayer(player);
                } else {
                    this.addOtherPlayer(player);
                }
            });
        });

        this.socket.on('disconnect', (playerInfo) => {
            playerInfo = JSON.parse(playerInfo);
            let playerId = playerInfo.id;

            this.otherPlayers.getChildren().forEach((otherPlayer) => {
                if (playerId === otherPlayer.playerId) {
                    otherPlayer.destroy();
                }
            });
        });

        this.socket.on('error', () => {
            console.log(arguments);
        });
    }

    update() {}

    createAnimations() {
        // Пока не используется

        let skins = ['lavaDark', 'lava', 'desert'];
        let skinsCount = skins.length;

        for (let i = 0; i < skinsCount; i++) {

            this.anims.create({
                key: skins[i] + '-move',
                frames: this.anims.generateFrameNumbers('tank', { start: 1 + 5 * i, end: 4 + 5 * i }),
                frameRate: 12,
                repeat: -1
            });
        }
    }

    addPlayer(playerInfo) {
        this.tank = this.physics.add.sprite(playerInfo.x, playerInfo.y, 'tank', 5).setOrigin(0.5, 0.5);
    }

    addOtherPlayer(playerInfo) {
        /**
         * TODO Событие newPlayer рассылается всем клиентам, включая тех, кто только подключился
         * Из-за этого танк игрока сначала рисуется в методе addPlayer (по событию currentPlayer)
         * А затем на его место еще рисуется как будто бы танк другого игрока (по событию newPlayer)
         * Пока добавил проверку, но по хорошему событие newPlayer должно рассылаться остальным клиентам исключая
         * подключившегося игрока
         */
        if (this.socket.id === playerInfo.id) {
            return;
        }

        const otherPlayer = this.add.sprite(playerInfo.x, playerInfo.y, 'tank', 10).setOrigin(0.5, 0.5);
        otherPlayer.playerId = playerInfo.id;
        this.otherPlayers.add(otherPlayer);
    }
}
