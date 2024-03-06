import React, { useState, useEffect } from 'react';

const BackendEndpoints = () => {
  const [expressions, setExpressions] = useState([]);
  
  const fetchExpressions = async () => {
    const response = await fetch('/api/go/expressions');
    const data = await response.json();
    setExpressions(data);
  };

  useEffect(() => {
    fetchExpressions();
  }, []);

  return (
    <div>
      <h1>Список выражений:</h1>
      <ul>
        {expressions.map(expression => (
          <li key={expression.id}>{expression.text}</li>
        ))}
      </ul>
    </div>
  );
};

export default BackendEndpoints;
