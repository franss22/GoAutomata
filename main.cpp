// #include <cstdlib>
// #include <iomanip>
// #include <iostream>
// #include <list>
// #include <vector>
// #include <map>
// #include <tuple>
// #include <math.h>
// #include <iostream>
// #include <unordered_map>

// using namespace std;

// #define NO_CC -1

// struct Word
// {
//     char letter;
//     int num;
// };

// struct ClockCondition
// {
//     int c1;
//     int c2;
// };

// struct ClockVals
// {
//     int z1;
//     int z2;
// };

// struct ClockResets
// {
//     bool reset_z1;
//     bool reset_z2;
// };

// struct Transition
// {
//     int p;
//     int q;
//     char l;
//     int c_1;
//     int c_2;
//     bool reset_z1;
//     bool reset_z2;
// };

// bool le(ClockVals &c1, ClockVals &c2)
// {
//     return c1.z1 <= c2.z1 && c1.z2 <= c2.z2;
// }

// bool is_minimal(std::vector<ClockVals> &vals, int index)
// {
//     int cur_idx = 0;
//     for (std::vector<ClockVals>::iterator it = vals.begin(); it != vals.end(); ++it)
//     {
//         if (cur_idx == index)
//         {
//             continue;
//         }
//         if (le(vals[cur_idx], vals[index]))
//         {
//             return false;
//         }
//         cur_idx++;
//     }
//     return true;
// }

// void get_pareto(std::vector<ClockVals> &vals,
//                 std::vector<ClockVals> &pareto)
// {
//     for (long unsigned int i = 0; i < vals.size(); i++)
//     {
//         if (is_minimal(vals, i))
//         {
//             pareto.push_back(vals[i]);
//         }
//     }
// }
// bool matches_clock_condition(ClockVals &cv, ClockCondition &cc)
// {
//     bool result = true;
//     if (cc.c1 != NO_CC)
//     {
//         result = result && cv.z1 <= cc.c1;
//     }
//     if (cc.c2 != NO_CC)
//     {
//         result = result && cv.z2 <= cc.c2;
//     }
//     return result;
// }

// class TimeCEA
// {
// public:
//     int Q;
//     unordered_map<pair<int, char>, vector<tuple<ClockCondition, ClockResets, int>>> Delta;
//     int q_0;
//     unordered_map<int, vector<ClockVals>> onthefly_states;
//     int Z;
//     TimeCEA(int Q)
//     {
//         this->Q = Q;
//         this->q_0 = 0;
//         this->Z = 2;
//     }

//     void add_transition(int p, char l, int c_1, int c_2, bool reset_z1, bool reset_z2, int q)
//     {
//         if (this->Delta.find(make_pair(p, l)) == this->Delta.end())
//         {
//             vector<tuple<ClockCondition, ClockResets, int>> trans;

//             this->Delta.insert(make_pair(make_pair(p, l), trans));
//         }
//         this->Delta[make_pair(p, l)].push_back(make_tuple((ClockCondition){c_1, c_2}, (ClockResets){reset_z1, reset_z2}, q));
//     }

//     void receive_word(Word word)
//     {
//         unordered_map<int, vector<ClockVals>> new_onthefly_states;

//         for (auto it = this->onthefly_states.begin(); it != this->onthefly_states.end(); it++)
//         {
//             int p = it->first;
//             for (auto it2 = it->second.begin(); it2 != it->second.end(); it2++)
//             {
//                 ClockVals clockval = *it2;
//                 execute_transition(p, clockval, word, new_onthefly_states);
//             }
//         }
//         new_onthefly_states[this->q_0].push_back((ClockVals){0, 0});
//         this->onthefly_states = new_onthefly_states;
//         for (auto it = this->onthefly_states.begin(); it != this->onthefly_states.end(); it++)
//         {
//             vector<ClockVals> pareto;
//             get_pareto(it->second, pareto);
//             it->second = pareto;
//         }
//     }

//     void execute_transition(int p, ClockVals &clocks, Word &word, unordered_map<int, vector<ClockVals>> &new_onthefly_states)
//     {
//         if (this->Delta.find(make_pair(p, word.letter)) != this->Delta.end())
//         {
//             for (auto it = this->Delta[make_pair(p, word.letter)].begin(); it != this->Delta[make_tuple(p, word.letter)].end(); it++)
//             {
//                 ClockVals new_cc{clocks.z1 + word.num, clocks.z2 + word.num};
//                 ClockCondition cc = get<0>(*it);
//                 ClockResets cr = get<1>(*it);
//                 int q = get<2>(*it);

//                 if (matches_clock_condition(new_cc, cc))
//                 {
//                     if (cr.reset_z1)
//                     {
//                         new_cc.z1 = 0;
//                     }
//                     if (cr.reset_z2)
//                     {
//                         new_cc.z2 = 0;
//                     }
//                     new_onthefly_states[q].push_back(new_cc);
//                 }
//             }
//         }
//     }
// };

// void generate_all_transitions(vector<Transition> &result, int paretonum, int Q)
// {
//     int max_clock_conditions = max(2 * Q, 2 * paretonum);
//     for (int c1 = 0; c1 <= max_clock_conditions; c1++)
//     {
//         for (int c2 = 0; c2 <= c1; c2++)
//         {
//             for (int reset_z1 = 0; reset_z1 <= 1; reset_z1++)
//             {

//                 for (int reset_z2 = 0; reset_z2 <= 1; reset_z2++)
//                 {
//                     char l = 'a';
//                     for (int p = 0; p < Q; p++)
//                     {
//                         for (int q = p; q < Q; q++)
//                         {
//                             result.push_back((Transition){p, q, l, c1 == 0 ? NO_CC : c1, c2 == 0 ? NO_CC : c2, !!reset_z1, !!reset_z2});
//                         }
//                     }
//                 }
//             }
//         }
//     }
// }

// int test_automata_for_pareto(TimeCEA &simple, int max_wlen, int paretonum)
// {
//     for (int length = 1; length <= max_wlen; length++)
//     {
//         simple.receive_word((Word){'a', 1});
//         for (auto it = simple.onthefly_states.begin(); it != simple.onthefly_states.end(); it++)
//         {
//             if (it->second.size() >= (long unsigned int)paretonum)
//             {
//                 return length + 1;
//             }
//         }
//     }

//     return -1;
// }

// int powerSet(vector<Transition> &transitions, int index, vector<Transition> &curr, int Q, int paretonum, int max_wlen)
// {
//     int n = transitions.size();
//     if (index == n - 1)
//     {
//         TimeCEA simple = TimeCEA(Q);
//         for (auto it = curr.begin(); it != curr.end(); it++)
//         {
//             simple.add_transition(it->p, it->l, it->c_1, it->c_2, it->reset_z1, it->reset_z2, it->q);
//         }
//         int length = test_automata_for_pareto(simple, max_wlen, paretonum);
//         if (length != -1)
//         {
//             return length;
//         }
//         return -1;
//     }
//     for (int i = index + 1; i < n; i++)
//     {

//         curr.push_back(transitions[i]);
//         int length = powerSet(transitions, i, curr, Q, paretonum, max_wlen);

//         if (length == -1)
//         {
//             curr.pop_back();
//         }
//         else
//         {
//             return length;
//         }
//     }
//     return -1;
// }

// ostream &operator<<(ostream &os, const TimeCEA &tcea)
// {
//     os << "nstates: " << tcea.Q << endl;
//     for (auto it = tcea.Delta.begin(); it != tcea.Delta.end(); it++)
//     {
//         int p = get<0>(it->first);
//         char l = get<1>(it->first);
//         os << "(" << p << ", " << l << ")" << " =>" << endl;
//         for (auto it2 = it->second.begin(); it2 != it->second.end(); it2++)
//         {
//             ClockCondition cc = get<0>(*it2);
//             ClockResets cr = get<1>(*it2);
//             int q = get<2>(*it2);

//             os << "    " << "(";
//             os << "z1 < " << cc.c1 << " and z2 < " << cc.c2;
//             os << ", ";
//             os << "r1: " << cr.reset_z1;
//             os << ", ";
//             os << "r2: " << cr.reset_z2;
//             os << ", ";
//             os << q;
//             os << ")" << endl;
//         }
//     }
//     exit(1);
//     return os;
// }

// int get_n_states_for_paretonum(int paretonum)
// {
//     for (int Q = 1; Q <= 1000; Q++)
//     {
//         cout << "Testing with " << Q << " states..." << endl;
//         int max_wlen = max(2 * Q, 2 * paretonum);
//         vector<Transition> transitions;
//         generate_all_transitions(transitions, paretonum, Q);
//         bool found = false;
//         vector<Transition> curr;
//         powerSet(transitions, -1, curr, Q, paretonum, max_wlen, found);
//         if (found)
//         {
//             TimeCEA simple = TimeCEA(Q);
//             for (auto it = curr.begin(); it != curr.end(); it++)
//             {
//                 simple.add_transition(it->p, it->l, it->c_1, it->c_2, it->reset_z1, it->reset_z2, it->q);
//             }
//             cout << simple;
//             return Q;
//         }
//     }
//     return -1;
// }

// int main(int argc, char *argv[])
// {
//     // vector<Transition> transitions;
//     // generate_all_transitions(transitions, 2, 2);

//     // for (auto it2 = transitions.begin(); it2 != transitions.end(); it2++)
//     // {
//     //     cout << "(" << it2->p << ", " << it2->l << ", " << "(" << it2->c_1 << ", " << it2->c_2 << ")" << ", " << "(" << it2->reset_z1 << ", " << it2->reset_z2 << ")" << ", " << it2->q << ")" << endl;
//     // }
//     // exit(0);
//     for (int paretonum = 1; paretonum <= 10; paretonum++)
//     {
//         cout << "Testing paretonum " << paretonum << endl;
//         int nstates = get_n_states_for_paretonum(paretonum);
//         if (nstates > 0)
//         {
//             cout << "paretonum:" << paretonum << endl;
//             cout << "N states: " << nstates << endl;
//         }
//     }
// }
