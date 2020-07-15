"use strict";

class TestScene extends Phaser.Scene
{
    _otherPlayerIds = [];

    preload() {
        this.load.spritesheet('tank', 'img/tank.png', { frameWidth: 60, frameHeight: 60 });
    }

    create() {
        this.socket = io();
        this.cursors = this.input.keyboard.createCursorKeys();
        this.otherPlayers = this.physics.add.group();
        this.createAnimations();

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

    update() {
        if (!this.tank) {
            return;
        }

        if (this.cursors.left.isDown || this.cursors.right.isDown || this.cursors.up.isDown || this.cursors.down.isDown) {
            let tankSpeed = 2;

            if (this.cursors.left.isDown) {
                this.tank.x -= tankSpeed;
                this.tank.angle = 270;
            } else if (this.cursors.right.isDown) {
                this.tank.x += tankSpeed;
                this.tank.angle = 90;
            } else if (this.cursors.up.isDown) {
                this.tank.y -= tankSpeed;
                this.tank.angle = 0;
            } else if (this.cursors.down.isDown) {
                this.tank.y += tankSpeed;
                this.tank.angle = 180;
            }

            this.tank.anims.play('lava-move', true);
            this.socket.emit('playerMovement', JSON.stringify({
                x: this.tank.x,
                y: this.tank.y,
                direction: this.tank.angle / 90
            }));
        } else {
            this.tank.anims.play('lava-idle', true);
        }

        this.physics.world.wrap(this.tank, 5);
    }

    createAnimations() {
        let skins = ['lavaDark', 'lava', 'desert'];
        let skinsCount = skins.length;
        let imagesCount = 8;

        for (let i = 0; i < skinsCount; i++) {
            this.anims.create({
                key: skins[i] + '-move',
                frames: this.anims.generateFrameNumbers('tank', { start: 1 + imagesCount * i, end: 4 + imagesCount * i }),
                frameRate: 12,
                repeat: -1
            });

            this.anims.create({
                key: skins[i] + '-idle',
                frames: this.anims.generateFrameNumbers('tank', { start: 5 + imagesCount * i, end: 7 + imagesCount * i }),
                frameRate: 8,
                repeat: -1
            });
        }
    }

    tankCollision() {
        console.log('collision');
    }

    addPlayer(playerInfo) {
        if (undefined !== this.tank) {
            this.tank.x = playerInfo.x;
            this.tank.y = playerInfo.y;
        } else {
            this.tank = this.physics.add.sprite(playerInfo.x, playerInfo.y, 'tank').setOrigin(0.5, 0.5);
            this.physics.add.overlap(this.tank, this.otherPlayers, this.tankCollision, null, this);
        }
    }

    addOtherPlayer(playerInfo) {
        if (!this._otherPlayerIds.includes(playerInfo.id)) {
            const otherPlayer = this.add.sprite(playerInfo.x, playerInfo.y, 'tank', 16).setOrigin(0.5, 0.5);
            otherPlayer.playerId = playerInfo.id;
            otherPlayer.anims.play('desert-idle', true);

            this.otherPlayers.add(otherPlayer);
            this._otherPlayerIds.push(playerInfo.id);
        } else {
            this.otherPlayers.getChildren().forEach((otherPlayer) => {
                if (playerInfo.id === otherPlayer.playerId) {
                    if ((otherPlayer.x !== playerInfo.x) || (otherPlayer.y !== playerInfo.y)) {
                        otherPlayer.anims.play('desert-move', true);
                        clearTimeout(otherPlayer.idleTimer);
                        otherPlayer.idleTimer = setTimeout(function(otherPlayer) {
                            otherPlayer.anims.play('desert-idle', true);
                        }, 50, otherPlayer);
                    }

                    otherPlayer.x = playerInfo.x;
                    otherPlayer.y = playerInfo.y;
                }
            });
        }
    }
}
