/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useState } from 'react';
import { API } from '../../helpers';
import { showError, showSuccess } from '../../helpers';

const TokenQuery = () => {
  const [tokenKey, setTokenKey] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);

  const handleQuery = async () => {
    if (!tokenKey.trim()) {
      showError('请输入密钥');
      return;
    }

    setLoading(true);
    setResult(null);

    try {
      const res = await API.get(`/api/query/token?key=${encodeURIComponent(tokenKey.trim())}`);
      const { success, message, data } = res.data;
      
      if (success) {
        setResult(data);
        showSuccess('查询成功');
      } else {
        showError(message || '查询失败');
      }
    } catch (error) {
      showError('查询失败，请检查网络连接');
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleQuery();
    }
  };

  const formatQuota = (quota) => {
    if (!quota && quota !== 0) return '0';
    return (quota / 500000).toFixed(2);
  };

  const formatDate = (timestamp) => {
    if (!timestamp) return '永不过期';
    const date = new Date(timestamp);
    return date.toLocaleString('zh-CN');
  };

  const getStatusText = (status) => {
    const statusMap = {
      1: '正常',
      2: '已禁用',
      3: '已过期',
      4: '已耗尽'
    };
    return statusMap[status] || '未知';
  };

  const getStatusColor = (status) => {
    const colorMap = {
      1: 'text-green-600',
      2: 'text-gray-600',
      3: 'text-red-600',
      4: 'text-orange-600'
    };
    return colorMap[status] || 'text-gray-600';
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        {/* 标题 */}
        <div className="text-center mb-8">
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-2">
            密钥额度查询
          </h1>
          <p className="text-gray-600 text-sm sm:text-base">
            输入您的密钥查询剩余额度和使用情况
          </p>
        </div>

        {/* 查询表单 */}
        <div className="bg-white rounded-2xl shadow-xl p-6 sm:p-8 mb-6">
          <div className="space-y-4">
            <div>
              <label htmlFor="token-key" className="block text-sm font-medium text-gray-700 mb-2">
                API 密钥
              </label>
              <input
                id="token-key"
                type="text"
                value={tokenKey}
                onChange={(e) => setTokenKey(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="请输入密钥，例如：sk-xxxxx"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all outline-none text-sm sm:text-base"
                disabled={loading}
              />
            </div>
            <button
              onClick={handleQuery}
              disabled={loading}
              className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-6 rounded-lg transition-colors duration-200 disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center justify-center space-x-2"
            >
              {loading ? (
                <>
                  <svg className="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <span>查询中...</span>
                </>
              ) : (
                <span>查询</span>
              )}
            </button>
          </div>
        </div>

        {/* 查询结果 */}
        {result && (
          <div className="bg-white rounded-2xl shadow-xl p-6 sm:p-8 animate-fadeIn">
            <h2 className="text-xl sm:text-2xl font-bold text-gray-900 mb-6 flex items-center">
              <svg className="w-6 h-6 mr-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              查询结果
            </h2>

            <div className="space-y-4">
              {/* 状态 */}
              <div className="flex justify-between items-center py-3 border-b border-gray-100">
                <span className="text-gray-600 font-medium">状态</span>
                <span className={`font-semibold ${getStatusColor(result.status)}`}>
                  {getStatusText(result.status)}
                </span>
              </div>

              {/* 总额度 */}
              <div className="flex justify-between items-center py-3 border-b border-gray-100">
                <span className="text-gray-600 font-medium">总额度</span>
                <span className="text-gray-900 font-semibold">
                  {result.unlimited_quota ? '无限' : `$${formatQuota(result.total_quota)}`}
                </span>
              </div>

              {/* 已使用 */}
              <div className="flex justify-between items-center py-3 border-b border-gray-100">
                <span className="text-gray-600 font-medium">已使用</span>
                <span className="text-red-600 font-semibold">
                  ${formatQuota(result.used_quota)}
                </span>
              </div>

              {/* 剩余额度 */}
              <div className="flex justify-between items-center py-3 border-b border-gray-100">
                <span className="text-gray-600 font-medium">剩余额度</span>
                <span className="text-green-600 font-semibold text-lg">
                  {result.unlimited_quota ? '无限' : `$${formatQuota(result.remain_quota)}`}
                </span>
              </div>

              {/* 过期时间 */}
              <div className="flex justify-between items-center py-3">
                <span className="text-gray-600 font-medium">过期时间</span>
                <span className="text-gray-900 font-semibold text-sm sm:text-base">
                  {formatDate(result.expired_time)}
                </span>
              </div>

              {/* 进度条（仅在非无限额度时显示） */}
              {!result.unlimited_quota && (
                <div className="mt-6">
                  <div className="flex justify-between text-sm text-gray-600 mb-2">
                    <span>使用进度</span>
                    <span>{((result.used_quota / result.total_quota) * 100).toFixed(1)}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
                    <div
                      className="h-full bg-gradient-to-r from-blue-500 to-purple-600 rounded-full transition-all duration-500"
                      style={{ width: `${Math.min((result.used_quota / result.total_quota) * 100, 100)}%` }}
                    ></div>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {/* 提示信息 */}
        <div className="mt-8 text-center text-sm text-gray-500">
          <p>💡 此查询功能无需登录，您可以随时查看密钥使用情况</p>
        </div>
      </div>

      <style jsx>{`
        @keyframes fadeIn {
          from {
            opacity: 0;
            transform: translateY(10px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        .animate-fadeIn {
          animation: fadeIn 0.3s ease-in-out;
        }
      `}</style>
    </div>
  );
};

export default TokenQuery;