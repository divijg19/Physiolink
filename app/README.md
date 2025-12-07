# PhysioLink Frontend (Flutter)

This project has been refactored from React Native to Flutter.

## Setup

1.  **Initialize Platform Code**:
    Since this project was generated without platform folders (android, ios, etc.), run:
    ```bash
    flutter create .
    ```

2.  **Install Dependencies**:
    ```bash
    flutter pub get
    ```

3.  **Generate Code**:
    Run the build runner to generate JSON serialization and Riverpod code:
    ```bash
    dart run build_runner build --delete-conflicting-outputs
    ```

## Architecture

*   **State Management**: Riverpod
*   **Navigation**: GoRouter
*   **HTTP Client**: Dio
*   **Storage**: Flutter Secure Storage
*   **Models**: Freezed & JSON Serializable

## Folder Structure

*   `lib/core`: Core infrastructure (API, Router, Theme).
*   `lib/features`: Feature-based modules (Auth, Home, etc.).
