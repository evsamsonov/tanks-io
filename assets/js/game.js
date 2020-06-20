const config = {
    type: Phaser.AUTO,
    width: 800,
    height: 600,
    physics: {
        default: 'arcade',
        arcade: {
            debug: false
        }
    },
    backgroundColor: '#bbbbbb',
    scene: [TestScene],
    title: 'Tanks',
    version: '0.03',
};
const game = new Phaser.Game(config);