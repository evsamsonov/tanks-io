const config = {
    type: Phaser.AUTO,
    width: 800,
    height: 600,
    physics: {
        default: 'arcade',
        arcade: {
            debug: false,
            gravity: { y: 0 }
        }
    },
    backgroundColor: '#bbbbbb',
    scene: [TestScene],
    title: 'Tanks',
    version: '0.02',
};
const game = new Phaser.Game(config);