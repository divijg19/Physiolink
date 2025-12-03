// src/utils/date.js
export const formatDateLong = (iso) => {
    if (!iso) return '';
    const d = new Date(iso);
    return d.toLocaleDateString(undefined, { weekday: 'long', month: 'long', day: 'numeric' });
};

export const formatTime = (iso) => {
    if (!iso) return '';
    const d = new Date(iso);
    return d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
};

export const toISODate = (iso) => {
    if (!iso) return '';
    return new Date(iso).toISOString().slice(0, 10);
};
