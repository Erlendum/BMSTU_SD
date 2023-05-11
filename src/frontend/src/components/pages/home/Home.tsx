import { FC, useState, useEffect } from 'react'

import styles from './Home.module.scss'
import { InstrumentService } from '../../../services/instrument.service'
import { useQuery } from 'react-query'
import InstrumentItem from '../../ui/instrument-item/InstrumentItem'
import Layout from '../../ui/layout/Layout'
import ReactPaginate from 'react-paginate'
import Modal from '../../ui/modal/Modal'
import AddInstrumentDBForm from '../../ui/form/AddInstrumentDBForm'
import { UserService } from '../../../services/user.service'
import { mergeAlias } from 'vite'

const Home: FC = () => {
	let isAdmin =
		UserService.getCurrentIsAdmin() != null &&
		UserService.getCurrentIsAdmin() == 'true'
	const [AddInstrumentInDBModalActive, setAddInstrumentInDBModalActive] =
		useState(false)

	const [updateQuery, setUpdateQuery] =
		useState(false)

	const { data: instruments, isLoading } = useQuery(
		['instruments', updateQuery],
		() => InstrumentService.getList(),
		{
			select: ({ instruments }) => instruments
		}
	)

	const [pageNumber, setPageNumber] = useState(0)
	const instrumentsPerPage = 24
	const pagesVisited = pageNumber * instrumentsPerPage

	const displayInstruments = instruments
		?.slice(pagesVisited, pagesVisited + instrumentsPerPage)
		.map(instrument => {
			return (
				<div>
					<InstrumentItem
						instrument={instrument}
						isComparisonList={false}
						key={instrument.InstrumentId}
						updateQuery={updateQuery}
						setUpdateQuery={setUpdateQuery}
					/>
				</div>
			)
		})

	// @ts-ignore
	const pagesCount = Math.ceil(instruments?.length / instrumentsPerPage)

	// @ts-ignore
	const changePage = ({ selected }) => {
		setPageNumber(selected)
	}

	return (
		<Layout title='Store collection'>
			{isLoading ? (
				<div className='text-black text-2xl'>Loading...</div>
			) : instruments?.length ? (
				<div>
					<button
						className={styles.openBtn}
						onClick={() => setAddInstrumentInDBModalActive(true)}
						hidden={!isAdmin}
					>
						Add instrument in DB
					</button>
					<div className={styles.wrapper}>{displayInstruments}</div>
					<ReactPaginate
						previousLabel={'Prev'}
						nextLabel={'Next'}
						pageCount={pagesCount}
						onPageChange={changePage}
						containerClassName={styles.paginationBttns}
						previousClassName={'previousBttns'}
						nextLinkClassName={'nextBttn'}
						disabledClassName={'paginationDisable'}
						activeClassName={styles.paginationActive}
					/>
					<Modal
						active={AddInstrumentInDBModalActive}
						setActive={setAddInstrumentInDBModalActive}
					>
						<AddInstrumentDBForm updateQuery={updateQuery} setUpdateQuery={setUpdateQuery}/>
					</Modal>
				</div>
			) : (
				<div>Instruments not found!</div>
			)}
		</Layout>
	)
}

export default Home
