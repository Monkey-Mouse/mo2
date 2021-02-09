export function randomProperty(obj: any) {
    const keys = Object.keys(obj);
    return obj[keys[keys.length * Math.random() << 0]];
}

export function Copy(mainObject: any) {
    const objectCopy = {}; // objectCopy will store a copy of the mainObject
    let key;
    for (key in mainObject) {
        objectCopy[key] = mainObject[key]; // copies each property to the objectCopy object
    }
    return objectCopy;
}