import Symbol
interface Data {
    prop: string[],
    value: number,
    type: symbol
}
let sym = Symbol();
let obj: Data = {
    prop: [],
    value: 0.123,
    type: Symbol('f')
};