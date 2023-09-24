const { DataTypes } = require('sequelize');

module.exports = (sequelize) => {
    return sequelize.define('Wallet', {
        WalletId: {
            type: DataTypes.STRING,
            primaryKey: true,
        },
        CreatedAt: {
            type: DataTypes.DATE,
            defaultValue: DataTypes.NOW,
        },
    }, {
        tableName: 'Wallet',
        timestamps: false, // Puoi impostare questo su true se desideri timestamps
    });
};
