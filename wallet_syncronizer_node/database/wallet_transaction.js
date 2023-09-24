const { DataTypes } = require('sequelize');
const Wallet = require('./wallet');

module.exports = (sequelize) => {
    const WalletTransaction = sequelize.define('WalletTransaction', {
        TxType: {
            type: DataTypes.STRING,
        },
        TxHash: {
            type: DataTypes.STRING,
        },
        Price: {
            type: DataTypes.FLOAT,
            allowNull: false,
        },
        Amount: {
            type: DataTypes.FLOAT,
        },
        Total: {
            type: DataTypes.FLOAT,
            allowNull: false,
        },
        AgeTimestamp: {
            type: DataTypes.DATE,
            allowNull: false,
        },
        Asset: {
            type: DataTypes.STRING,
            allowNull: false,
        },
        WalletId: {
            type: DataTypes.STRING,
        },
        CreatedAt: {
            type: DataTypes.DATE,
            defaultValue: DataTypes.NOW,
        },
    }, {
        tableName: 'WalletTransaction',
        timestamps: false, // Puoi impostare questo su true se desideri timestamps
    });

    // Definisci la relazione con il modello Wallet
    WalletTransaction.belongsTo(Wallet(sequelize), {
        foreignKey: 'WalletId',
        as: 'Wallet',
    });

    return WalletTransaction;
};
