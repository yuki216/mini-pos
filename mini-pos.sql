/*
 Navicat Premium Data Transfer

 Source Server         : DB Training
 Source Server Type    : MySQL
 Source Server Version : 50736
 Source Host           : 159.65.138.79:3306
 Source Schema         : polkesban

 Target Server Type    : MySQL
 Target Server Version : 50736
 File Encoding         : 65001

 Date: 14/11/2021 23:45:45
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `description` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `create_at` datetime(0) NOT NULL,
  `update_at` datetime(0) NOT NULL,
  `outlet_id` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `categories_fk0`(`outlet_id`) USING BTREE,
  CONSTRAINT `categories_fk0` FOREIGN KEY (`outlet_id`) REFERENCES `outlets` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of categories
-- ----------------------------
INSERT INTO `categories` VALUES (1, 'Elecronik 1', 'Elektronik Arus Lemah', '2021-11-13 18:30:48', '2021-11-13 18:30:52', 1);
INSERT INTO `categories` VALUES (3, 'Elektronik 1', 'Elektronik Arus Lemah', '2021-11-13 18:33:16', '2021-11-13 18:33:20', 2);

-- ----------------------------
-- Table structure for customers
-- ----------------------------
DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(75) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `email` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `phone` varchar(15) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `address` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `create_at` datetime(0) NOT NULL,
  `update_at` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of customers
-- ----------------------------
INSERT INTO `customers` VALUES (1, 'Yuki achriansyah', 'fachry.yuki@gmail.com', '6285314481882', 'Jl. Kopo Gg. Panineungan 1 No 207', '2021-11-14 07:51:54', '2021-11-14 07:55:53');
INSERT INTO `customers` VALUES (2, 'Mimi Jamilah', 'mimi@gmail.com', '6281392483988', 'Jl. Kopo Gg. Panineungan 1 No 207', '2021-11-14 07:53:30', '2021-11-14 07:53:30');

-- ----------------------------
-- Table structure for orderItem
-- ----------------------------
DROP TABLE IF EXISTS `orderItem`;
CREATE TABLE `orderItem`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `qty` int(11) NOT NULL,
  `discount` int(11) NULL DEFAULT NULL,
  `tax_apply` int(11) NULL DEFAULT NULL,
  `create_at` datetime(0) NOT NULL,
  `update_at` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `orderItem_fk0`(`order_id`) USING BTREE,
  INDEX `orderItem_fk1`(`product_id`) USING BTREE,
  CONSTRAINT `orderItem_fk0` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `orderItem_fk1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of orderItem
-- ----------------------------
INSERT INTO `orderItem` VALUES (1, 1, 1, 7, NULL, NULL, '2021-11-14 16:20:04', '2021-11-14 18:10:26');
INSERT INTO `orderItem` VALUES (3, 13, 1, 2, NULL, NULL, '2021-11-14 18:59:28', '2021-11-14 19:00:41');
INSERT INTO `orderItem` VALUES (5, 14, 1, 2, NULL, NULL, '2021-11-14 19:09:45', '2021-11-14 19:09:51');
INSERT INTO `orderItem` VALUES (6, 14, 2, 1, NULL, NULL, '2021-11-14 19:09:58', '2021-11-14 19:09:58');

-- ----------------------------
-- Table structure for orders
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NULL DEFAULT NULL,
  `isCheckout` tinyint(1) NULL DEFAULT 0,
  `create_at` datetime(0) NOT NULL,
  `update_at` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk1_orders`(`customer_id`) USING BTREE,
  CONSTRAINT `fk1_orders` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE = InnoDB AUTO_INCREMENT = 18 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of orders
-- ----------------------------
INSERT INTO `orders` VALUES (1, 1, 1, '2021-11-14 16:19:24', '2021-11-14 19:02:12');
INSERT INTO `orders` VALUES (13, 1, 1, '2021-11-14 18:52:14', '2021-11-14 18:52:14');
INSERT INTO `orders` VALUES (14, 1, 1, '2021-11-14 19:09:45', '2021-11-14 19:10:11');

-- ----------------------------
-- Table structure for outlet_products
-- ----------------------------
DROP TABLE IF EXISTS `outlet_products`;
CREATE TABLE `outlet_products`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NOT NULL,
  `sku` varchar(25) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `price` decimal(12, 2) NOT NULL,
  `quantity` float(255, 0) NOT NULL,
  `quantity_use` float(255, 0) NULL DEFAULT 0,
  `outlet_id` int(11) NOT NULL,
  `supplier_id` int(11) NOT NULL,
  `create_at` datetime(0) NULL DEFAULT NULL,
  `update_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `sku_unique`(`sku`) USING BTREE,
  INDEX `op_fk1`(`product_id`) USING BTREE,
  INDEX `op_fk2`(`outlet_id`) USING BTREE,
  INDEX `op_fk3`(`supplier_id`) USING BTREE,
  CONSTRAINT `op_fk1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `op_fk2` FOREIGN KEY (`outlet_id`) REFERENCES `outlets` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `op_fk3` FOREIGN KEY (`supplier_id`) REFERENCES `suppliers` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of outlet_products
-- ----------------------------
INSERT INTO `outlet_products` VALUES (1, 1, 'AMB1DS12301', 50000.00, 5, 5, 1, 1, '2021-11-14 12:48:55', '2021-11-14 13:52:17');
INSERT INTO `outlet_products` VALUES (2, 2, 'AMB1D12302', 65000.00, 5, 0, 1, 1, '2021-11-14 13:14:15', '2021-11-14 13:14:17');

-- ----------------------------
-- Table structure for outlets
-- ----------------------------
DROP TABLE IF EXISTS `outlets`;
CREATE TABLE `outlets`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `address` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `create_at` datetime(0) NOT NULL,
  `update_at` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `address`(`address`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of outlets
-- ----------------------------
INSERT INTO `outlets` VALUES (1, 'Outlet 1', 'Jl. ABCD berarti', '2021-11-13 18:29:02', '2021-11-13 18:29:09');
INSERT INTO `outlets` VALUES (2, 'Outlet 2', 'Jl. CSDSD kjdhdfsd', '2021-11-13 18:29:32', '2021-11-13 18:29:36');

-- ----------------------------
-- Table structure for payments
-- ----------------------------
DROP TABLE IF EXISTS `payments`;
CREATE TABLE `payments`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(11) NOT NULL,
  `tax` float NULL DEFAULT NULL,
  `type_payment` varchar(15) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `total_payment` decimal(12, 2) NOT NULL,
  `create_at` datetime(0) NOT NULL,
  `update_at` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `payments_fk0`(`order_id`) USING BTREE,
  CONSTRAINT `payments_fk0` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of payments
-- ----------------------------
INSERT INTO `payments` VALUES (1, 1, 0, 'cash', 50000.00, '2021-11-14 17:52:09', '2021-11-14 17:52:09');
INSERT INTO `payments` VALUES (2, 1, 0, 'cash', 50000.00, '2021-11-14 19:02:12', '2021-11-14 19:02:12');
INSERT INTO `payments` VALUES (3, 14, 0, 'cash', 50000.00, '2021-11-14 19:10:11', '2021-11-14 19:10:11');

-- ----------------------------
-- Table structure for products
-- ----------------------------
DROP TABLE IF EXISTS `products`;
CREATE TABLE `products`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `description` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `color` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `size` varchar(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `unit` varchar(25) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `image` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `category_id` int(11) NOT NULL,
  `create_at` datetime(0) NOT NULL,
  `update_at` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `products_fk0`(`category_id`) USING BTREE,
  CONSTRAINT `products_fk0` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of products
-- ----------------------------
INSERT INTO `products` VALUES (1, 'Samsung A3', '-', 'white', '-', 'EA', '/sadasdsa.jpg', 1, '2021-11-13 18:59:04', '2021-11-13 18:59:07');
INSERT INTO `products` VALUES (2, 'Samsung A7', '-', 'black', '-', 'EA', '/sadasdsa.jpg', 1, '2021-11-13 18:59:04', '2021-11-13 18:59:07');
INSERT INTO `products` VALUES (4, 'Samsung A5', '-', 'white', '-', 'EA', '/sadasdsa.jpg', 1, '2021-11-14 01:11:00', '2021-11-14 11:41:45');

-- ----------------------------
-- Table structure for purchases
-- ----------------------------
DROP TABLE IF EXISTS `purchases`;
CREATE TABLE `purchases`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NOT NULL,
  `cost_price` decimal(12, 2) NOT NULL,
  `quantity` float(255, 0) NOT NULL,
  `total_cost_price` decimal(12, 2) NOT NULL,
  `discount` float(255, 0) NULL DEFAULT 0,
  `tax` float(255, 0) NULL DEFAULT 5,
  `category_id` int(11) NOT NULL,
  `supplier_id` int(11) NOT NULL,
  `outlet_id` int(11) NOT NULL,
  `create_at` datetime(0) NULL DEFAULT NULL,
  `update_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `purchases_fk1`(`product_id`) USING BTREE,
  INDEX `purchases_fk2`(`category_id`) USING BTREE,
  INDEX `purchases_fk3`(`supplier_id`) USING BTREE,
  INDEX `purchases_fk4`(`outlet_id`) USING BTREE,
  CONSTRAINT `purchases_fk1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `purchases_fk2` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `purchases_fk3` FOREIGN KEY (`supplier_id`) REFERENCES `suppliers` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `purchases_fk4` FOREIGN KEY (`outlet_id`) REFERENCES `outlets` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of purchases
-- ----------------------------
INSERT INTO `purchases` VALUES (1, 1, 50000.00, 10, 505000.00, 0, 5, 1, 1, 1, '2021-11-14 10:51:00', '2021-11-14 11:08:02');
INSERT INTO `purchases` VALUES (2, 1, 50000.00, 10, 505000.00, 0, 5, 1, 1, 1, '2021-11-14 11:01:14', '2021-11-14 11:01:14');

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `id` int(11) NOT NULL,
  `name` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, 'super admin');
INSERT INTO `roles` VALUES (2, 'merchant');

-- ----------------------------
-- Table structure for suppliers
-- ----------------------------
DROP TABLE IF EXISTS `suppliers`;
CREATE TABLE `suppliers`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `address` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `code` varchar(5) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `create_at` datetime(0) NOT NULL,
  `update_at` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of suppliers
-- ----------------------------
INSERT INTO `suppliers` VALUES (1, 'PD. Ambang Raya 1', 'Jl. Kopo Gg. Panineungan 1 No 207', 'AMB', '2021-11-13 23:53:15', '2021-11-14 07:54:48');
INSERT INTO `suppliers` VALUES (2, 'PD. Ambang Raya 2', 'Jl asdsadas', 'AMB', '2021-11-13 23:53:15', '2021-11-13 23:53:19');

-- ----------------------------
-- Table structure for user_outlets
-- ----------------------------
DROP TABLE IF EXISTS `user_outlets`;
CREATE TABLE `user_outlets`  (
  `user_id` int(11) NOT NULL,
  `outlet_id` int(11) NOT NULL,
  INDEX `fk2`(`outlet_id`) USING BTREE,
  INDEX `fk1`(`user_id`) USING BTREE,
  CONSTRAINT `fk2` FOREIGN KEY (`outlet_id`) REFERENCES `outlets` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_outlets
-- ----------------------------
INSERT INTO `user_outlets` VALUES (7, 1);
INSERT INTO `user_outlets` VALUES (7, 2);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fullname` varchar(150) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `username` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `email` varchar(75) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `role_id` int(11) NOT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username`) USING BTREE,
  UNIQUE INDEX `email`(`email`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8 COLLATE = utf8_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (7, 'Admin', 'fachry45', 'admin@admin.com', '$2a$10$uybUyIVafxKlzBLi/1juSe91sauKjUBjUzRkTt361CNusWJMzF1Xq', 1, NULL, NULL);

SET FOREIGN_KEY_CHECKS = 1;
