"use strict";

class TestScene extends Phaser.Scene
{
    _otherPlayerIds = [];

    preload() {
        this.load.spritesheet('tank', 'img/tank.png', { frameWidth: 60, frameHeight: 60 });
    }

    create() {
        this.socket = io();
        this.otherPlayers = this.physics.add.group();

        this.socket.on('currentPlayers', (players) => {
            players = JSON.parse(players);

            players.forEach((player) => {
                if (player.id === this.socket.id ) {
                    this.addPlayer(player);
                } else {
                    this.addOtherPlayer(player);
                }
            });
        });

        this.socket.on('disconnect', (playerInfo) => {
            playerInfo = JSON.parse(playerInfo);
            let playerId = playerInfo.id;
            let playerIndex = this._otherPlayerIds.indexOf(playerId);

            if (playerIndex !== -1) {
                this.otherPlayers.getChildren().forEach((otherPlayer) => {
                    if (playerId === otherPlayer.playerId) {
                        otherPlayer.destroy();
                    }
                });
                this._otherPlayerIds.splice(playerIndex, 1);
            }
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
        if (undefined !== this.tank) {
            this.tank.x = playerInfo.x;
            this.tank.y = playerInfo.y;
        } else {
            this.tank = this.physics.add.sprite(playerInfo.x, playerInfo.y, 'tank', 5).setOrigin(0.5, 0.5);
        }
    }

    addOtherPlayer(playerInfo) {
        if (!this._otherPlayerIds.includes(playerInfo.id)) {
            const otherPlayer = this.add.sprite(playerInfo.x, playerInfo.y, 'tank', 10).setOrigin(0.5, 0.5);
            otherPlayer.playerId = playerInfo.id;

            this.otherPlayers.add(otherPlayer);
            this._otherPlayerIds.push(playerInfo.id);
        }
    }
}
