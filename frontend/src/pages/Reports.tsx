import { useState, useEffect } from 'react';
import Layout from '../components/Layout';
import { transactionApi, Transaction } from '../api/transactions';
import { storeApi, Store } from '../api/stores';

const Reports = () => {
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [stores, setStores] = useState<Store[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [dateRange, setDateRange] = useState('all');
    const [selectedStore, setSelectedStore] = useState<number | 'all'>('all');

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        setLoading(true);
        setError('');
        try {
            const [txData, storeData] = await Promise.all([
                transactionApi.getTransactions(1, 100),
                storeApi.getStores(),
            ]);
            setTransactions(txData.transactions || []);
            setStores(storeData);
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to load data');
        } finally {
            setLoading(false);
        }
    };

    // Filter transactions
    const filteredTransactions = transactions.filter(tx => {
        if (selectedStore !== 'all' && tx.store_id !== selectedStore) {
            return false;
        }
        return true;
    });

    // Calculate summary statistics
    const totalIncrease = filteredTransactions
        .filter(tx => tx.transaction_type === 'INCREASE')
        .reduce((sum, tx) => sum + tx.total_amount, 0);
    const totalDecrease = filteredTransactions
        .filter(tx => tx.transaction_type === 'DECREASE')
        .reduce((sum, tx) => sum + tx.total_amount, 0);

    const summary = {
        totalTransactions: filteredTransactions.length,
        totalIncrease,
        totalDecrease,
        increaseCount: filteredTransactions.filter(tx => tx.transaction_type === 'INCREASE').length,
        decreaseCount: filteredTransactions.filter(tx => tx.transaction_type === 'DECREASE').length,
        totalRevenue: totalDecrease - totalIncrease,
    };

    // Group transactions by store
    const salesByStore = stores.map(store => {
        const storeTx = filteredTransactions.filter(
            tx => tx.transaction_type === 'DECREASE' && tx.store_id === store.id
        );
        return {
            store: store.name,
            transactions: storeTx.length,
            totalSales: storeTx.reduce((sum, tx) => sum + tx.total_amount, 0),
        };
    }).filter(s => s.transactions > 0);

    // Group transactions by product
    const productStats = filteredTransactions.reduce((acc, tx) => {
        const key = tx.product_name || `Product ${tx.product_id}`;
        if (!acc[key]) {
            acc[key] = {
                name: key,
                increase: 0,
                decrease: 0,
                revenue: 0,
            };
        }
        if (tx.transaction_type === 'INCREASE') {
            acc[key].increase += tx.quantity;
        } else {
            acc[key].decrease += tx.quantity;
            acc[key].revenue += tx.total_amount;
        }
        return acc;
    }, {} as Record<string, { name: string; increase: number; decrease: number; revenue: number }>);

    const topProducts = Object.values(productStats)
        .sort((a, b) => b.revenue - a.revenue)
        .slice(0, 5);

    return (
        <Layout>
            <div className="px-4 py-6 sm:px-0">
                <div className="flex justify-between items-center mb-6">
                    <h1 className="text-3xl font-bold text-gray-900">Transaction Reports</h1>
                    <button
                        onClick={loadData}
                        className="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700 transition-colors"
                    >
                        Refresh
                    </button>
                </div>

                {/* Filters */}
                <div className="bg-white shadow rounded-lg p-4 mb-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Store Filter
                            </label>
                            <select
                                value={selectedStore}
                                onChange={(e) => setSelectedStore(e.target.value === 'all' ? 'all' : parseInt(e.target.value))}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                            >
                                <option value="all">All Stores</option>
                                {stores.map(store => (
                                    <option key={store.id} value={store.id}>{store.name}</option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Date Range
                            </label>
                            <select
                                value={dateRange}
                                onChange={(e) => setDateRange(e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                            >
                                <option value="all">All Time</option>
                                <option value="today">Today</option>
                                <option value="week">This Week</option>
                                <option value="month">This Month</option>
                            </select>
                        </div>
                    </div>
                </div>

                {/* Error */}
                {error && (
                    <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
                        {error}
                    </div>
                )}

                {/* Loading */}
                {loading ? (
                    <div className="text-center py-12">
                        <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
                        <p className="mt-4 text-gray-600">Loading reports...</p>
                    </div>
                ) : (
                    <>
                        {/* Summary Cards */}
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
                            <div className="bg-white shadow rounded-lg p-6">
                                <div className="flex items-center">
                                    <div className="flex-shrink-0 bg-blue-500 rounded-md p-3">
                                        <svg className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                                        </svg>
                                    </div>
                                    <div className="ml-4">
                                        <p className="text-sm font-medium text-gray-500">Total Transactions</p>
                                        <p className="text-2xl font-bold text-gray-900">{summary.totalTransactions}</p>
                                    </div>
                                </div>
                            </div>

                            <div className="bg-white shadow rounded-lg p-6">
                                <div className="flex items-center">
                                    <div className="flex-shrink-0 bg-green-500 rounded-md p-3">
                                        <svg className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 11l5-5m0 0l5 5m-5-5v12" />
                                        </svg>
                                    </div>
                                    <div className="ml-4">
                                        <p className="text-sm font-medium text-gray-500">Stock Purchases</p>
                                        <p className="text-2xl font-bold text-gray-900">฿{summary.totalIncrease.toLocaleString()}</p>
                                        <p className="text-xs text-gray-500">{summary.increaseCount} transactions</p>
                                    </div>
                                </div>
                            </div>

                            <div className="bg-white shadow rounded-lg p-6">
                                <div className="flex items-center">
                                    <div className="flex-shrink-0 bg-orange-500 rounded-md p-3">
                                        <svg className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 13l-5 5m0 0l-5-5m5 5V6" />
                                        </svg>
                                    </div>
                                    <div className="ml-4">
                                        <p className="text-sm font-medium text-gray-500">Store Sales</p>
                                        <p className="text-2xl font-bold text-gray-900">฿{summary.totalDecrease.toLocaleString()}</p>
                                        <p className="text-xs text-gray-500">{summary.decreaseCount} transactions</p>
                                    </div>
                                </div>
                            </div>

                            <div className="bg-white shadow rounded-lg p-6">
                                <div className="flex items-center">
                                    <div className="flex-shrink-0 bg-purple-500 rounded-md p-3">
                                        <svg className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                    </div>
                                    <div className="ml-4">
                                        <p className="text-sm font-medium text-gray-500">Gross Profit</p>
                                        <p className="text-2xl font-bold text-gray-900">฿{summary.totalRevenue.toLocaleString()}</p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
                            {/* Sales by Store */}
                            <div className="bg-white shadow rounded-lg p-6">
                                <h3 className="text-lg font-semibold text-gray-900 mb-4">Sales by Store</h3>
                                {salesByStore.length > 0 ? (
                                    <div className="space-y-4">
                                        {salesByStore.map((item, index) => {
                                            const maxSales = Math.max(...salesByStore.map(s => s.totalSales));
                                            const percentage = (item.totalSales / maxSales) * 100;
                                            return (
                                                <div key={index}>
                                                    <div className="flex justify-between text-sm mb-1">
                                                        <span className="font-medium text-gray-700">{item.store}</span>
                                                        <span className="text-gray-900 font-semibold">฿{item.totalSales.toLocaleString()}</span>
                                                    </div>
                                                    <div className="flex items-center">
                                                        <div className="w-full bg-gray-200 rounded-full h-2 mr-2">
                                                            <div
                                                                className="bg-primary-600 h-2 rounded-full transition-all duration-300"
                                                                style={{ width: `${percentage}%` }}
                                                            ></div>
                                                        </div>
                                                        <span className="text-xs text-gray-500 whitespace-nowrap">{item.transactions} tx</span>
                                                    </div>
                                                </div>
                                            );
                                        })}
                                    </div>
                                ) : (
                                    <p className="text-gray-500 text-center py-4">No sales data available</p>
                                )}
                            </div>

                            {/* Top Products */}
                            <div className="bg-white shadow rounded-lg p-6">
                                <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Selling Products</h3>
                                {topProducts.length > 0 ? (
                                    <div className="space-y-4">
                                        {topProducts.map((product, index) => (
                                            <div key={index} className="flex items-center justify-between pb-3 border-b last:border-b-0">
                                                <div className="flex items-center">
                                                    <div className="flex-shrink-0 w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
                                                        <span className="text-primary-600 font-semibold text-sm">{index + 1}</span>
                                                    </div>
                                                    <div className="ml-3">
                                                        <p className="text-sm font-medium text-gray-900">{product.name}</p>
                                                        <p className="text-xs text-gray-500">{product.decrease} units sold</p>
                                                    </div>
                                                </div>
                                                <div className="text-right">
                                                    <p className="text-sm font-semibold text-gray-900">฿{product.revenue.toLocaleString()}</p>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                ) : (
                                    <p className="text-gray-500 text-center py-4">No product data available</p>
                                )}
                            </div>
                        </div>

                        {/* Transaction History Table */}
                        <div className="bg-white shadow rounded-lg overflow-hidden">
                            <div className="px-6 py-4 border-b border-gray-200">
                                <h3 className="text-lg font-semibold text-gray-900">Recent Transactions</h3>
                            </div>
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Date
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Type
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Product
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Store
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Quantity
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Amount
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {filteredTransactions.slice(0, 20).map((tx) => (
                                            <tr key={tx.id} className="hover:bg-gray-50">
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                                    {(() => {
                                                        const date = new Date(tx.transaction_date);
                                                        // Subtract 7 hours to convert from UTC to Bangkok time
                                                        date.setHours(date.getHours() - 7);
                                                        return date.toLocaleString('th-TH', {
                                                            year: 'numeric',
                                                            month: 'short',
                                                            day: 'numeric',
                                                            hour: '2-digit',
                                                            minute: '2-digit',
                                                            second: '2-digit'
                                                        });
                                                    })()}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap">
                                                    <span className={`px-2 py-1 text-xs font-semibold rounded-full ${tx.transaction_type === 'INCREASE'
                                                        ? 'bg-green-100 text-green-800'
                                                        : 'bg-orange-100 text-orange-800'
                                                        }`}>
                                                        {tx.transaction_type}
                                                    </span>
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                                    {tx.product_name || `Product ${tx.product_id}`}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                                    {tx.store_name || '-'}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                                    {tx.quantity}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                                    ฿{tx.total_amount.toLocaleString()}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </Layout>
    );
};

export default Reports;
