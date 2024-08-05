    --table: users
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    degree VARCHAR(255),
    grad_year VARCHAR(4),
    current_job VARCHAR(255),
    phone VARCHAR(20),
    session_key VARCHAR(255) UNIQUE,
    profile_picture VARCHAR(255),
    linkedin_profile VARCHAR(255),
    twitter_profile VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

--table: portfolio order table
CREATE TABLE portfolio_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    security VARCHAR(255) NOT NULL,
    unit INT NOT NULL,
    status VARCHAR(255),
    cancelled INT DEFAULT 0,
    user_id INT,
    user_email VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE
);


--table: transactions
CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_user_id INT,
    from_user_email VARCHAR(255),
    to_user_id INT,
    to_user_email VARCHAR(255),
    type VARCHAR(255),
    amount INT,
    user_email VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user_id) REFERENCES users(id),
    FOREIGN KEY (from_user_email) REFERENCES users(email),
    FOREIGN KEY (to_user_id) REFERENCES users(id),
    FOREIGN KEY (to_user_email) REFERENCES users(email)
);

--table: forums
CREATE TABLE forums (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    author VARCHAR(255),
    slug VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (author) REFERENCES users(email) ON DELETE SET NULL
);

--table: chat_messages
CREATE TABLE chat_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender INT,
    recipient INT,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (sender) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (recipient) REFERENCES users(id) ON DELETE SET NULL
);
