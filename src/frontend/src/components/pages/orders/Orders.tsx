import React, { FC, useEffect, useState } from 'react'
import { useQuery } from 'react-query'
import { UserService } from '../../../services/user.service'
import { toast, ToastContainer } from 'react-toastify'
import { OrderService } from '../../../services/order.service'
import InstrumentItemTable from '../../ui/instrument-item-table/InstrumentItemTable'
import Layout from '../../ui/layout/Layout'
import styles from '../orders/Orders.module.scss'
import OrderItem from '../../ui/order-item/OrderItem'
import OrderElementItem from '../../ui/orderElement-item/OrderElementItem'
import orderElementItem from '../../ui/orderElement-item/OrderElementItem'
import { IOrder } from '../../../types/order.interface'

const Orders: FC = () => {
	let isAdmin =
		UserService.getCurrentIsAdmin() != null &&
		UserService.getCurrentIsAdmin() == 'true'
	const [updateQuery, setUpdateQuery] =
		useState(false)
	const [orderId, setOrderId] = useState(-1)


	const { data: orders, isLoading } = useQuery(
		['orders', updateQuery],
		() => !isAdmin ? OrderService.getList() : OrderService.getListForAll(),
		{
			select: ({ orders }) => orders
		}
	)

	const displayOrders = orders?.map(
		(order: IOrder) => {
			return (
				<OrderItem
					updateQuery={updateQuery}
					setUpdateQuery={setUpdateQuery}
					order={order}
					key={order.OrderId}
				/>)

		}
	)

	return (
		<Layout title='Orders'>
			<table className={styles.table}>
				<thead>
				<tr>
					<th className={styles.textLeft}>Order Id</th>
					<th className={styles.textLeft}>Date</th>
					<th className={styles.textLeft}>Price</th>
					<th className={styles.textLeft}>Status</th>
					<th hidden={!isAdmin} className={styles.textLeft}>User Id</th>
					<th className={styles.textLeft}>
						Links
					</th>
				</tr>
				</thead>
				<tbody className='table-hover'>
				{displayOrders}
				</tbody>
			</table>
			<ToastContainer />
		</Layout>
	)
}

export default Orders
