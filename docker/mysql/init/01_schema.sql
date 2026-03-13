CREATE DATABASE IF NOT EXISTS visual_lab;
USE visual_lab;

CREATE TABLE IF NOT EXISTS employees (
  id INT PRIMARY KEY,
  dept_id INT NOT NULL,
  age INT NOT NULL,
  name VARCHAR(64) NOT NULL,
  salary DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_dept_age_name (dept_id, age, name)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS engine_demo_myisam (
  id INT PRIMARY KEY,
  note VARCHAR(64)
) ENGINE=MyISAM;

CREATE TABLE IF NOT EXISTS engine_demo_memory (
  id INT PRIMARY KEY,
  note VARCHAR(64)
) ENGINE=MEMORY;

INSERT INTO employees (id, dept_id, age, name, salary) VALUES
  (1, 10, 22, 'Alice', 12000.00),
  (2, 10, 31, 'Bob', 15000.00),
  (3, 10, 35, 'Carol', 18200.00),
  (4, 20, 25, 'David', 13400.00),
  (5, 20, 41, 'Eva', 19800.00),
  (6, 30, 29, 'Frank', 14100.00)
ON DUPLICATE KEY UPDATE
  dept_id = VALUES(dept_id),
  age = VALUES(age),
  name = VALUES(name),
  salary = VALUES(salary);
