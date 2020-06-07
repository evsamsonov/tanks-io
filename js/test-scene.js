"use strict";

class TestScene extends Phaser.Scene
{
    preload() {
        this.load.spritesheet('tank', 'assets/tank.png', { frameWidth: 60, frameHeight: 60 });
    }

    create() {
        let skins = ['lavaDark', 'lava', 'desert'];
        let skinsCount = skins.length;

        for (let i = 0; i < skinsCount; i++) {
            let animationPrefix = skins[i] + '-';

            this.anims.create({
                key: animationPrefix + 'move',
                frames: this.anims.generateFrameNumbers('tank', { start: 1 + 5 * i, end: 4 + 5 * i }),
                frameRate: 12,
                repeat: -1
            });
        }

        this.tanks = [
            this.add.sprite(0, 0, 'tank').setOrigin(0.5, 0.5).setAngle(90).anims.play('lava-move', true),
            this.add.sprite(0, 0, 'tank').setOrigin(0.5, 0.5).setAngle(90).anims.play('desert-move', true),
        ];
    }

    update() {
        const state = getCurrentState();

        state.players.forEach((item, i) => {
            this.tanks[i].setPosition(item.x, item.y);
        });
    }
}
