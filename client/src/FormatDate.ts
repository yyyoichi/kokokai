export const getYYYYMMDD = (d: Date, separator: string = "") => {
    const date = new FormatDate(d);
    const format = [date.getFullYear(), date.getMonth0(), date.getDate0()];
    return format.join(separator);
};
export default class FormatDate {
    constructor(public d: Date) {}
    getMonth0() {
        const d = "0" + String(this.d.getMonth() + 1);
        return d.slice(-2);
    }
    getDate0() {
        const d = "0" + String(this.d.getDate());
        return d.slice(-2);
    }
    getFullYear() {
        return this.d.getFullYear();
    }
}
