const config = {
    type: Phaser.AUTO,
    width: 800,
    height: 600,
    backgroundColor: '#bbbbbb',
    scene: [TestScene],
    title: 'Tanks',
    version: '0.01',
};
const game = new Phaser.Game(config);