const { DataTypes } = require('sequelize');

module.exports = (sequelize) => {
    return sequelize.define('Token', {
        TokenId: {
            type: DataTypes.STRING,
            primaryKey: true,
        },
        Name: {
            type: DataTypes.STRING,
            allowNull: false,
        },
        Symbol: {
            type: DataTypes.STRING,
            allowNull: false,
        },
        Decimals: {
            type: DataTypes.INTEGER,
            allowNull: false,
        },
        CreatedAt: {
            type: DataTypes.DATE,
            defaultValue: DataTypes.NOW,
        },
        Logo: {
            type: DataTypes.STRING,
        },
        GoPlusResponse: {
            type: DataTypes.BLOB, // Assuming it's binary data
        },
    }, {
        tableName: 'Token',
        timestamps: false, // You can set this to true if you want timestamps
    });
};