import React, { useState } from 'react';
import './Home.css';

export const Home = () => {
    const [orderNumber, setOrderNumber] = useState('');
    const [orderData, setOrderData] = useState(null);
    const [error, setError] = useState(null);

    const handleInputChange = (e) => {
    setOrderNumber(e.target.value);
    };

    const handleAddOrder = async () => {
    try {
        const response = await fetch(`/api/orders/${orderNumber}`);
        const data = await response.json();
        setOrderData(data);
        setOrderNumber('');
    } catch (err) {
        setError('Ошибка при получении данных о заказе');
    }
    };

    return (
    <div className="home">
        <div>
        <label htmlFor="orderNumber">Номер заказа:</label>
            <input
            type="text"
            id="orderNumber"
            name="orderNumber"
            placeholder="Введите номер заказа"
            value={orderNumber}
            onChange={handleInputChange}
            />
        <button onClick={handleAddOrder}>Добавить заказ</button>
        </div>
            {error && <p className="error">{error}</p>}
            {orderData && (
                <div className="order-details">
                    <h3>Данные о заказе:</h3>
                    <ul>
            {Object.entries(orderData).map(([key, value]) => (
                <li key={key}>
                <strong>{key}:</strong> {value}
                </li>
            ))}
            </ul>
        </div>
        )}
    </div>
    );
};
