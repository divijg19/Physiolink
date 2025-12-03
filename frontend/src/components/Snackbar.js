import React, { useEffect, useRef } from 'react';
import { View, Text, StyleSheet, Animated } from 'react-native';
import { COLORS, SPACING, FONT } from '../theme';

const Snackbar = ({ message, visible = false, duration = 3000 }) => {
    const opacity = useRef(new Animated.Value(0)).current;

    useEffect(() => {
        let timer;
        if (visible) {
            Animated.timing(opacity, { toValue: 1, duration: 180, useNativeDriver: true }).start();
            timer = setTimeout(() => {
                Animated.timing(opacity, { toValue: 0, duration: 180, useNativeDriver: true }).start();
            }, duration);
        }
        return () => clearTimeout(timer);
    }, [visible]);

    if (!visible) return null;

    return (
        <Animated.View style={[styles.container, { opacity }]} pointerEvents="none">
            <Text style={styles.text}>{message}</Text>
        </Animated.View>
    );
};

const styles = StyleSheet.create({
    container: {
        position: 'absolute',
        bottom: SPACING.lg,
        left: SPACING.md,
        right: SPACING.md,
        backgroundColor: COLORS.darkGray || '#333',
        padding: SPACING.sm,
        borderRadius: 10,
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 999,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.12,
        shadowRadius: 12,
        elevation: 6,
    },
    text: {
        color: '#fff',
        fontFamily: 'Poppins_600SemiBold',
        fontSize: FONT.body,
    },
});

export default Snackbar;
