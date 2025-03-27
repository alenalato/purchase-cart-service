import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:9090/api';

export default function () {
    let validOrderPayload = JSON.stringify({
        order: {
            items: [
                { product_id: 1, quantity: 1 },
                { product_id: 2, quantity: 5 },
                { product_id: 3, quantity: 1 }
            ]
        }
    });

    let invalidOrderPayload422 = JSON.stringify({
        order: {
            items: [
                { product_id: 5, quantity: 1 }
            ]
        }
    });

    let invalidOrderPayload400 = JSON.stringify({
        order: {
            items: [
                { product_id: 1, quantity: -1 }
            ]
        }
    });

    let params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Test valid order
    let res = http.post(`${BASE_URL}/v1/order`, validOrderPayload, params);

    check(res, {
        'is status 201': (r) => r.status === 201,
        'response contains order id': (r) => r.json().hasOwnProperty('id'),
        'response contains total price': (r) => r.json().hasOwnProperty('total_price') && r.json().total_price === 12.5,
        'response contains total vat': (r) => r.json().hasOwnProperty('total_vat') && r.json().total_vat === 1.25,
        'response contains items': (r) => r.json().hasOwnProperty('items'),
    });

    if (res.json().hasOwnProperty('items')) {
        let index = 0;
        let item = res.json().items[index];
        check(item, {
            [`item ${index} contains product_id`]: (i) => i.hasOwnProperty('product_id') && i.product_id === 1,
            [`item ${index} contains quantity`]: (i) => i.hasOwnProperty('quantity') && i.quantity === 1,
            [`item ${index} contains price`]: (i) => i.hasOwnProperty('price') && i.price === 2,
            [`item ${index} contains vat`]: (i) => i.hasOwnProperty('vat') && i.vat === 0.2,
        });

        index = 1;
        item = res.json().items[index];
        check(item, {
            [`item ${index} contains product_id`]: (i) => i.hasOwnProperty('product_id') && i.product_id === 2,
            [`item ${index} contains quantity`]: (i) => i.hasOwnProperty('quantity') && i.quantity === 5,
            [`item ${index} contains price`]: (i) => i.hasOwnProperty('price') && i.price === 7.5,
            [`item ${index} contains vat`]: (i) => i.hasOwnProperty('vat') && i.vat === 0.75,
        });

        index = 2;
        item = res.json().items[index];
        check(item, {
            [`item ${index} contains product_id`]: (i) => i.hasOwnProperty('product_id') && i.product_id === 3,
            [`item ${index} contains quantity`]: (i) => i.hasOwnProperty('quantity') && i.quantity === 1,
            [`item ${index} contains price`]: (i) => i.hasOwnProperty('price') && i.price === 3,
            [`item ${index} contains vat`]: (i) => i.hasOwnProperty('vat') && i.vat === 0.3,
        });
    }

    // Test invalid order (422)
    let invalidRes422 = http.post(`${BASE_URL}/v1/order`, invalidOrderPayload422, params);

    check(invalidRes422, {
        'is status 422': (r) => r.status === 422,
        'response contains proper error code': (r) => r.json().code === 'unprocessable_entity',
        'response contains error message': (r) => r.json().hasOwnProperty('message'),
    });

    // Test invalid order (400)
    let invalidRes400 = http.post(`${BASE_URL}/v1/order`, invalidOrderPayload400, params);

    check(invalidRes400, {
        'is status 400': (r) => r.status === 400,
        'response contains proper error code': (r) => r.json().code === 'invalid_argument',
        'response contains error message': (r) => r.json().hasOwnProperty('message'),
    });

    sleep(1);
}
