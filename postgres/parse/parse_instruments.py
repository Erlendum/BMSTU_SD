import csv
import random
import time

import requests
import numpy as np
from bs4 import BeautifulSoup

def rand_material():
    materials = ['Липа', 'Дуб', 'Сосна', 'Пластик', 'Береза', 'Яблоня', 'Ольха', 'Махонь', 'Орех', 'Клен']
    i = random.randint(0, len(materials) - 1)
    return materials[i]

def append_one_elem(instrument, id, ids, names, prices, materials, types, brands, imgs):
    ids.append(id)
    try:
        img = instrument.find_all('img')[0]['data-src']
    except KeyError:
        img = 'https://novosibirsk.orso-b.ru/img/nophoto.jpg'
    try:
        price = int(instrument['data-price'])
    except KeyError:
        price = random.randint(2000, 100000)
    soup = instrument.find_all('a', href=True, itemprop=False, class_=False)

    types.append(soup[0].contents[0])
    names.append(soup[1].contents[0])
    brands.append(names[-1].split()[0])
    prices.append(price)
    imgs.append(img)
    materials.append(rand_material())

    return True

instruments_list = []

with open('instruments_links.txt') as f:
    instruments_links = f.readlines()

for i in range(len(instruments_links)):
    instruments_links[i] = instruments_links[i].split()[0]

for instruments_link in instruments_links:
    n = 1
    n_page = 1
    while n > 0:
        response = requests.get(instruments_link + f'?in-stock=1&pre-order=1&page={n_page}')
        soup = BeautifulSoup(response.text, 'lxml')
        instruments = soup.find_all('section', class_='product-thumbnail')

        for instrument in instruments:
            if len(instrument['class']) == 1:
                instruments_list.append(instrument)
        n = len(instruments)
        print(instruments_link + '\t' + f'page {n_page}' + '\r', end='')
        n_page += 1

ids = ['instrument_id']
names = ['instrument_name']
prices = ['instrument_price']
materials = ['instrument_material']
types = ['instrument_type']
brands = ['instrument_brand']
imgs = ['instrument_img']

id = 1

for instrument in instruments_list:
    if append_one_elem(instrument, id, ids, names, prices, materials, types, brands, imgs) is None:
        id -= 1
    print(str(id) + '/' + str(len(instruments_list)) + '\r', end='')
    id += 1

np.savetxt('instruments.csv', [p for p in zip(ids, names, prices, materials, types, brands, imgs)], delimiter=';', fmt='%s', encoding='utf8')
