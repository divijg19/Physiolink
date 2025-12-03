import React, { useState } from 'react';
import { View, Text, Image, StyleSheet } from 'react-native';
import { COLORS } from '../theme';

const Avatar = ({ uri, name, size = 72 }) => {
    const [errored, setErrored] = useState(false);

    const initials = (name || '')
        .split(' ')
        .map(s => s[0])
        .filter(Boolean)
        .slice(0, 2)
        .join('')
        .toUpperCase();

    if (uri && !errored) {
        return (
            <Image
                source={{ uri }}
                style={[styles.image, { width: size, height: size, borderRadius: size / 2 }]}
                onError={() => setErrored(true)}
                resizeMode='cover'
            />
        );
    }

    return (
        <View style={[styles.fallback, { width: size, height: size, borderRadius: size / 2, backgroundColor: COLORS.primary }]}>
            <Text style={[styles.initials, { fontSize: Math.round(size / 2.8) }]}>{initials || 'PT'}</Text>
        </View>
    );
};

const styles = StyleSheet.create({
    image: {
        backgroundColor: COLORS.lightGray,
    },
    fallback: {
        justifyContent: 'center',
        alignItems: 'center',
    },
    initials: {
        color: '#fff',
        fontFamily: 'Poppins_600SemiBold',
    },
});

export default Avatar;
