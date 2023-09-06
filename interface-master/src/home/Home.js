import React, { useState } from 'react';
import './Home.css';

export const Home = () => {
    const [orderUID, setOrderUID] = useState('');
    const [orderData, setOrderData] = useState(null);
    const [error, setError] = useState(null);

    const handleInputChange = (e) => {
        setOrderUID(e.target.value);
    };

    const handleAddOrder = async () => {
        try {
            const response = await fetch(`http://localhost:8000/orders/${orderUID}`);
            if (!response.ok) {
                throw new Error('Ошибка при получении данных о заказе');
            }
            const data = await response.json();
            setOrderData(data);
            setOrderUID('');
        } catch (err) {
            setError(err.message);
        }
    };

    return (
        <div className="home">
            <div>
                <label htmlFor="orderUID">Номер заказа:</label>
                <input
                    type="text"
                    id="orderUID"
                    name="orderUID"
                    placeholder="Введите номер заказа"
                    value={orderUID}
                    onChange={handleInputChange}
                />
                <button onClick={handleAddOrder}>Найти заказ</button>
            </div>
            {error && <p className="error">{error}</p>}
            {orderData && (
                <div className="order-details">
                    <h3>Данные о заказе:</h3>
                    <ul>
                        {Object.entries(orderData).map(([key, value]) => (
                            <li key={key}>
                                <strong>{key}:</strong> {JSON.stringify(value)}
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
};
